#!/bin/bash
curl 'https://www.cosmosgranite.com/getProductDetail' \
  -H 'X-requested-with: XMLHttpRequest' \
  --data-raw 'name=Titanium&location=charlotte&id=20488&pro_link=https%3A%2F%2Fwww.cosmosgranite.com%2Fcharlotte%2Fgranite%2Fcharlotte-293-titanium' \
  > titanium.charlotte.json
