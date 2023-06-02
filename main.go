package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
)

const json1 = `[{
            "uuid": "00000000010000000000000000000003",
            "name": "Default Policy",
            "showname": "_T_Default Policy",
            "enable": true,
            "__id": 0,
            "group": "00000000010000000000000000000001",
            "src": {
                "srcZones": ["12"],
                "srcAddrs": {
                    "srcAddrType": "NETOBJECT",
                    "srcIpGroups": [
                        "5FE8FDA10D8343CB9310E8B252E94C74",
                        "5FE8FDA10D8343CB9310E8B25274"
                    ]
                }
            },
            "dst": {
                "dstZones": ["12"],
                "dstAddrs": {
                    "dstAddrType": "NETOBJECT",
                    "dstIpGroups": [
                        "5FE8FDA10D8343CB9310E8B252E94C74"
                    ]
                },
                "services": [
                    "00000000000000000000000100000000"
                ],
                "applications": [
                    "00000000000000000000000600000001"
                ]
            },
            "action": 2,
            "schedule": "00000000000000000000000c00000001",
            "isdefault": true,
            "hits": 1,
"test":45,
"word":"tuyyuu"
        },{
            "uuid": "00000000010000000000000000000003",
            "name": "Default Policy",
            "showname": "_T_Default Policy",
            "enable": true,
            "__id": 0,
            "group": "00000000010000000000000400000001",
            "src": {
                "srcZones": ["21"],
                "srcAddrs": {
                    "srcAddrType": "NETOBJRSD",
                    "srcIpGroups": [
                        "5FE8FDA10D8343CB9310E8B25274"
                    ]
                }
            },
            "dst": {
                "dstZones": ["11"],
                "dstAddrs": {
                    "dstAddrType": "NEBJECT",
                    "dstIpGroups": [
                        "5FE8FDA10D8343CB9B252E94C74"
                    ]
                },
                "services": [
                    "00000000000000000000000100000000"
                ],
                "applications": [
                    "00000000000000000000000600000001"
                ],"application": [
                    "00000000000000000000000600000001"
                ]
            },
            "action": 1,
            "schedule": "00000000000000000000000c00000001",
            "isdefault": true,
            "hits": 2
        }]`

const json2 = `[{
            "uuid": "0000000001000000002300000000003",
            "name": "Default Policy",
            "showname": "_T_Default Policy",
            "enable": true,
            "__id": 0,
            "group": "00000000010000000000000400000001",
            "src": {
                "srcZones": ["21"],
                "srcAddrs": {
                    "srcAddrType": "NETOBJRSD",
                    "srcIpGroups": [
                        "5FE8FDA10D8343CB9310E8B25274"
                    ]
                }
            },
            "dst": {
                "dstZones": ["11"],
                "dstAddrs": {
                    "dstAddrType": "NEBJECT",
                    "dstIpGroups": [
                        "5FE8FDA10D8343CB9B252E94C74"
                    ]
                },
                "services": [
                    "00000000000000000000000100000000"
                ],
                "applications": [
                    "00000000000000000000000600000001"
                ],"application": [
                    "00000000000000000000000600000001"
                ]
            },
            "action": 1,
            "schedule": "00000000000000000000000c00000001",
            "isdefault": true,
            "hits": 2
        },{
            "uuid": "0000000001004500000000000000003",
            "name": "Default Policy",
            "showname": "_T_Default Policy",
            "enable": true,
            "__id": 0,
            "group": "00000000010000000000000400000001",
            "src": {
                "srcZones": ["21"],
                "srcAddrs": {
                    "srcAddrType": "NETOBJRSD",
                    "srcIpGroups": [
                        "5FE8FDA10D8343CB9310E8B25274"
                    ]
                }
            },
            "dst": {
                "dstZones": ["11"],
                "dstAddrs": {
                    "dstAddrType": "NEBJECT",
                    "dstIpGroups": [
                        "5FE8FDA10D8343CB9B252E94C74"
                    ]
                },
                "services": [
                    "00000000000000000000000100000000"
                ],
                "applications": [
                    "00000000000000000000000600000001"
                ],"application": [
                    "00000000000000000000000600000001"
                ]
            },
            "action": 1,
            "schedule": "00000000000000000000000c00000001",
            "isdefault": true,
            "hits": 2,
			"word":"tuyyuu"
        }]`

func main() {
	var obj1, obj2 map[string]interface{}
	_ = json.Unmarshal([]byte(json1), &obj1)
	_ = json.Unmarshal([]byte(json2), &obj2)

	/// 计算差异

	//diff := cmp.Diff(obj1, obj2, cmp.Ignore())
	//_ = cmp.Diff(obj1, obj2)

	//m := gjson.Parse(json1).Array()
	array := gjson.Parse(json1).IsArray()
	obj := gjson.Parse(json1).IsObject()
	value := gjson.Parse(json1).Type
	fmt.Println(array, obj, value)
	// 解析JSON字符串

}
