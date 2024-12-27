package agent

import (
    "golang.org/x/text/language"
    "golang.org/x/text/message"
)

var languages = []language.Tag{
    language.Japanese,
    language.English,
}

func InitCommandMessages() {
    languages := map[string]string{
        "Enter site key":
        "サイトキーを入力してくださ", 
        "Enter password":
        "アクセスキーを入力してください", 
        "js":
        "ジャバスクリプト",
    }
    for message_en, message_jp := range languages {
        _ = message.SetString(language.Japanese, message_en, message_jp)
    }
}

func Translate(acceptLanguage string, msg string, args ...interface{}) string {
    t, _, _ := language.ParseAcceptLanguage(acceptLanguage)
    matcher := language.NewMatcher(languages)
    tag, _, _ := matcher.Match(t...)
    p := message.NewPrinter(tag)
    return p.Sprintf(msg, args...)
}

