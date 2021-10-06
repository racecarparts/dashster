const calendarInterval = 3600000 // every hour

function writeCalendar(calData) {
    let text = calData.title + "<br>";
    text += calData.cal

    writeWidget('calendar', text)
}

function loadCalendar() {
    fetch('/calendar', {
        method: 'get'
    })
        .then(r => r.json())
        .then(jsonData => {
            writeCalendar(jsonData)
        })
        .catch(err => {
            console.log(err)
        })
}

loadCalendar()
setInterval(() => {
    loadCalendar()
}, calendarInterval)