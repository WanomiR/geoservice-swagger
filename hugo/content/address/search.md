---
menu:
    after:
        name: search
        weight: 0
title: Geo API
---

# Geo API
<link rel="stylesheet" href="https://unpkg.com/leaflet@1.7.1/dist/leaflet.css" crossorigin=""/>

<p>Поиск адреса</p>
<input id="search" />

<div id="result"></div>

<div id="mapid" style="height: 50vh"></div>

## Документация

Маршрут: `/api/address/search` метод `POST`
```go
type SearchRequest struct {
    Query string `json:"query"`
}
```

```go
type SearchResponse struct {
    Addresses []*Address `json:"addresses"`
}
```

Маршрут: `/api/address/geocode` метод `POST`
```go
type GeocodeRequest struct {
    Lat string `json:"lat"`
    Lng string `json:"lng"`
}
```

```go
type GeocodeResponse struct {
    Addresses []*Address `json:"addresses"`
}
```

## Провайдер
API: https://dadata.ru/api/ 

<!-- Include Leaflet JavaScript -->
<script src="https://unpkg.com/leaflet@1.7.1/dist/leaflet.js" crossorigin=""></script>
<script>
    let startPos = [59.9311, 30.3609];
    var mymap = L.map('mapid').setView(startPos, 11);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: 'Map data &copy; OpenStreetMap contributors',
        maxZoom: 18
    }).addTo(mymap);
    var currentMarker = null;
    // Обработчик события клика по карте
    mymap.on('click', function(e) {
        let data = {
            lat: e.latlng.lat.toString(),
            lng: e.latlng.lng.toString()
        };
        fetch('/api/address/geocode', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => response.json())
        .then(data => {
           table.setData(data.addresses);
           if (data.addresses.length > 0) {
                mymap.flyTo([data.addresses[0].lat, data.addresses[0].lon], 17);
                if (currentMarker) {
                    // Перемещение существующего маркера
                    currentMarker.setLatLng({lat: data.addresses[0].lat, lng: data.addresses[0].lon});
                } else {
                    // Создание нового маркера
                    currentMarker = L.marker({lat: data.addresses[0].lat, lng: data.addresses[0].lon}).addTo(mymap);
                }
           }
        })
        .catch(error => {
            console.log('Error:', error);
        });
    });
    // Сброс текущего маркера при двойном клике
    mymap.on('dblclick', function(e) {
        console.log('dblclick');
        if (currentMarker) {
            mymap.removeLayer(currentMarker);
            currentMarker = null;
        }
    });
</script>
<link href="https://unpkg.com/tabulator-tables@5.5.0/dist/css/tabulator.min.css" rel="stylesheet">
 <script type="text/javascript" src="https://unpkg.com/tabulator-tables@5.5.0/dist/js/tabulator.min.js"></script>
<script type="text/javascript">
//Build Tabulator
var tableData = [];
var table = new Tabulator("#result", {
    height:"311px",
    layout:"fitColumns",
    reactiveData:true, //turn on data reactivity
    responsiveLayout: "hide",
    data:tableData, //assign data to table
    placeholder:"No Data Set",
    selectable: true,
    autoColumns:true, //create columns from data field names
    rowClick:function(e, cell) {
        e.preventDefault();
        console.log("rowClick fired");
        e.stopPropagation();
    },
    selectableCheck:function(row){
        //row - row component
        let data = row.getData();
        if (data.lat != "" && data.lon != "") {
            mymap.flyTo([data.lat, data.lon], 17);
            if (currentMarker) {
                // Перемещение существующего маркера
                currentMarker.setLatLng({lat: data.lat, lng: data.lon});
            } else {
                // Создание нового маркера
                currentMarker = L.marker({lat: data.lat, lng: data.lon}).addTo(mymap);
            }
        }
        table.deselectRow();
        console.log("select fired");
        return true; //allow selection of rows where the age is greater than 18
    },
});
document.getElementById('search').addEventListener('input', function() {
    console.log('search change');
    if (this.value.length < 3) {
        return;
    }
    const data = {
        query: this.value
    };
    fetch('/api/address/search', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
       table.setData(data.addresses);
       if (data.addresses.length > 0) {
            mymap.flyTo([data.addresses[0].lat, data.addresses[0].lon], 17);
       }
    })
    .catch(error => {
        console.log('Error:', error);
    });
});
</script>
