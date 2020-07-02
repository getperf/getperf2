package agent

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Datastore struct {
	StartTime time.Time
	Host      string
	StatName  string
	DateDir   string
	TimeDir   string
	OutDir    string
}

type DatastoreKey struct {
	Host     string `json:"host"`
	StatName string `json:"stat_name"`
}

type DatastoreSet struct {
	DsKey   DatastoreKey `json:"keys"`
	ZipFile string
}

func NewDatastoreBase(outDir, host, statName string, start time.Time) *Datastore {
	datastore := &Datastore{
		StartTime: start,
		Host:      host,
		StatName:  statName,
		DateDir:   GetTimeString(YYYYMMDD, start),
		TimeDir:   GetTimeString(HHMISS, start),
		OutDir:    outDir,
	}
	return datastore
}

func (c *Config) NewDatastore(statName string, start time.Time) (*Datastore, error) {
	schedule := c.Schedule
	if schedule == nil {
		return nil, fmt.Errorf("new out log schedule not found")
	}
	collector := schedule.Collectors[statName]
	if collector == nil {
		return nil, fmt.Errorf("new out log collector not found %s", statName)
	}
	interval := time.Duration(collector.StatInterval)
	start = start.Truncate(time.Second * interval)
	host := c.GetServiceOrHostName()
	datastore := NewDatastoreBase(c.OutDir, host, statName, start)
	return datastore, nil
}

func (c *Config) NewDatastoreCurrent(statName string) (*Datastore, error) {
	return c.NewDatastore(statName, time.Now())
}

func (o *Datastore) RelDir() string {
	return filepath.Join(o.StatName, o.DateDir, o.TimeDir)
}

func (o *Datastore) AbsDir() string {
	return filepath.Join(o.OutDir, o.RelDir())
}

func (o *Datastore) ZipFile(host string) string {
	return fmt.Sprintf("arc_%s__%s_%s_%s.zip",
		host, o.StatName, o.DateDir, o.TimeDir,
	)
}

func (o *Datastore) ZipFilePrefix(host string) string {
	return fmt.Sprintf("arc_%s__%s_", host, o.StatName)
}

func (o *Datastore) OldZipFile(host string, hour int) string {
	start := o.StartTime.Add(-1 * time.Hour * time.Duration(hour))
	dateDir := GetTimeString(YYYYMMDD, start)
	timeDir := GetTimeString(HHMISS, start)
	return fmt.Sprintf("arc_%s__%s_%s_%s.zip", host, o.StatName, dateDir, timeDir)
}

func (ds *Datastore) GetZipFiles(since string) ([]string, error) {
	var files []string
	zipFilePrefix := ds.ZipFilePrefix(ds.Host)
	sinceKey := fmt.Sprintf("%s%s.zip", zipFilePrefix, since)
	zipFiles, err := ioutil.ReadDir(ds.OutDir)
	if err != nil {
		return files, errors.Wrap(err, "search zip files")
	}
	for _, zipFile := range zipFiles {
		zipFileName := zipFile.Name()
		if since != "" && strings.Compare(zipFileName, sinceKey) <= 0 {
			continue
		}
		if strings.HasPrefix(zipFileName, zipFilePrefix) {
			files = append(files, zipFileName)
		}
	}
	return files, nil
}

// OldestZipはアーカイブディレクトリから sinceラベルよりも後に生成された最も古い zip
// ファイルを検索し、そのファイル名を返します。since は YYYYMMDD_HHMM 形式で指定します

func (ds *Datastore) OldestZip(since string) (string, error) {
	files, err := ds.GetZipFiles(since)
	if err != nil {
		return "", errors.Wrap(err, "find oldest zip")
	}
	if len(files) == 0 {
		return "", errors.New("zip not found")
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	return files[0], nil
}

// LatestZipはアーカイブディレクトリから、最新の zip ファイルを検索し、そのファイル名を返します

func (ds *Datastore) LatestZip() (string, error) {
	files, err := ds.GetZipFiles("")
	if err != nil {
		return "", errors.Wrapf(err, "find latest zip")
	}
	if len(files) == 0 {
		return "", errors.New("zip not found")
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})

	return files[0], nil
}

func (ds *Datastore) GetDatastoreSets() (*[]DatastoreSet, error) {
	datastoreSets := []DatastoreSet{}
	zipFiles, err := ioutil.ReadDir(ds.OutDir)
	if err != nil {
		return &datastoreSets, errors.Wrap(err, "get keys of zip files")
	}
	dsKeysMap := make(map[DatastoreKey]string)
	for _, zipFile := range zipFiles {
		zipFileName := zipFile.Name()
		if !strings.HasPrefix(zipFileName, "arc_") ||
			!strings.HasSuffix(zipFileName, ".zip") {
			continue
		}
		// parse 'arc_{host}__{job}_{date}_{time}.zip'
		zipNames := strings.Split(zipFileName, "_")
		if len(zipNames) != 6 {
			continue
		}
		dsKey := DatastoreKey{
			Host:     zipNames[1],
			StatName: zipNames[3],
		}
		zipFileOld, ok := dsKeysMap[dsKey]
		if !ok {
			dsKeysMap[dsKey] = zipFileName
		} else {
			if strings.Compare(zipFileName, zipFileOld) > 0 {
				dsKeysMap[dsKey] = zipFileName
			}
		}
	}
	for dsKey, zipFile := range dsKeysMap {
		datastoreSet := DatastoreSet{dsKey, zipFile}
		datastoreSets = append(datastoreSets, datastoreSet)
	}
	return &datastoreSets, nil
}
