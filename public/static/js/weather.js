const weatherInterval = 1800000 // every 30 minutes
let weatherIntervalId;

function writeWeather(weatherData) {
    let weather = ""
    for (let i = 0; i < weatherData.length; i++) {
        let location = weatherData[i].location.padEnd(16) + "  "
        let cond = weatherData[i].cur_conditions.padEnd(16) + "  "
        let temp = weatherData[i].cur_temp.padEnd(5) + "\n"
        weather += location + cond + temp
    }

    writeWidget('weather-interval', timeIntervalStr(weatherInterval))
    writeWidget('weather', weather)
}

function loadWeather() {
    disableButton('load-weather')
    writeWidget('weather-interval', 'loading...')
    fetch('/weather', {
        method: 'get'
    })
        .then(r => r.json())
        .then(jsonData => {
            writeWeather(jsonData)
            enableButton('load-weather')
            weatherIntervalId = setupInterval(weatherIntervalId, weatherInterval, loadWeather)

        })
        .catch(err => {
            writeWidget('weather-interval', 'error: ' + err)
            enableButton('load-weather')
            console.log(err)
        })
}

loadWeather()
