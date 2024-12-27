import logging.config

from requests import Session
from zeep import Client
from zeep.transports import Transport

logging.config.dictConfig(
    {
        "version": 1,
        "formatters": {"verbose": {"format": "%(name)s: %(message)s"}},
        "handlers": {
            "console": {
                "level": "DEBUG",
                "class": "logging.StreamHandler",
                "formatter": "verbose",
            },
        },
        "loggers": {
            "zeep.transports": {
                "level": "DEBUG",
                "propagate": True,
                "handlers": ["console"],
            },
        },
    }
)
session = Session()
# session.cert = "./network/client.pem"
session.verify = "./network/ca.crt"
transport = Transport(session=session)
client = Client(
    "http://192.168.231.160:57000/axis2/services/GetperfService?wsdl",
    transport=transport,
)
result = client.service.registAgent(
    "site1", "host1", "e5c3a8cdc5db95b22517457e9750f5fc949eaaa3"
)
zipdata = result.attachments[0].content

# print(f"result: {result}\n")
# print(f"zip={zipdata}\n")

f = open('sslconf.zip', 'wb')
f.write(zipdata)
f.close()

