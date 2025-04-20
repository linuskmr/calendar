# Calendar

Get the timezone UTC offset for a specific datetime in a specific a timezone name:

```js
new Intl.DateTimeFormat("de-DE", {"timeZone": "Europe/Helsinki", "timeZoneName": "longOffset"}).format(new Date("2025-04-10T10:00:00+01:00"))
```


## UI

Features:

- Event list of the next 7 days
- Add new event



## API

### GET `/api/events?start={start}&end={end}&timezone={timezone}`

Returns a list of events for the given calendar name, between the start and end dates, in the given timezone.


### GET `/api/event/{id}`

Returns the event with the given id.


### POST `/api/event`

Creates a new event from the JSON data provided in the body:

```json
{
  "title": "",
  "start": "YYYY-MM-DDTHH:MM:SSZ",
  "end": "YYYY-MM-DDTHH:MM:SSZ",
  "description": "",
  "location": "",
  "calendar": ""
}
```

### PATCH `/api/event/{id}`

Updates the event with the specified id with the JSON data provided in the body.
