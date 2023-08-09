#!/bin/bash

# This has the metadata about all the slab types
curl 'https://ohm.stoneprofits.com/FetchDataWebV1.ashx?act=getItemGallery&InventoryGroupBy=IDTwo_&SearchbyItemIdentifiers=on&ShowFeatureProductOnTop=null&OnHold=null&OnSO=null&Intransit=null&showNotInStock=null&SearchbyFinish=on&SearchbySKU=on&Alphabet=' \
  | jq . > allSlabs.json

# This JSON has the details of individual lots of a given slab type
curl 'https://ohm.stoneprofits.com/FetchDataWebV1.ashx?act=getItemInventory&id=5181&InventoryGroupBy=IDTwo_&TrimmedUserID=4932186393528091&OnHold=null&OnSO=null&Intransit=null&SelectedLocation=&ShowLocationinGallery=on&LotPicturesRestrictToSIPL=False&ShowOnlyFullInventoryImages=on' \
  | jq . > copacabana.white.3cm.json
