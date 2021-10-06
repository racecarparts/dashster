const daylightMapInterval = 300000 // 5 minutes

function calcTimes(timeZones) {
    let d = new Date()
    let options = {
            // timeZone: 'Europe/London',
            // year: 'numeric',
            // month: 'numeric',
            // day: 'numeric',
            hour: 'numeric',
            minute: 'numeric',
            // second: 'numeric',
        }
    let options24Time = {
        hour: 'numeric',
        minute: 'numeric',
        hourCycle: 'h23',
        hour12: false
    }

    let times = []
    let group = 0
    for (let i = 0; i < timeZones.length; i++) {
        let timeZone = timeZones[i]
        if (timeZone.group !== group) {
            times.push("")
            group++
        }
        let time = ""
        if (timeZone.current_zone === true) {
            time += "*"
        } else {
            time += "&nbsp;"
        }
        options.timeZone = timeZone.tz
        options24Time.timeZone = timeZone.tz
        let formatter = new Intl.DateTimeFormat([], options);
        let formatter24 = new Intl.DateTimeFormat([], options24Time)

        let clockTime = formatter.format(d)
        if (clockTime.length === 7) {
            clockTime = "&nbsp;" + clockTime
        }
        let clockTime24 = formatter24.format(d)

        time += timeZone.short_tz + " " + clockTime + " (" + clockTime24 + ") " + timeZone.utc_offset
        times.push(time)
    }

    return times
}

function writeClock(data) {
    let text = ""
    let times = calcTimes(data)
    for (let i = 0; i < times.length; i++) {
        text += times[i] + "<br>";
    }
    writeWidget('world_clock', text)
}

function loadClock() {
    fetch('/worldclock', {
        method: 'get'
    })
        .then(r => r.json())
        .then(jsonData => {
            writeClock(jsonData)
            setInterval(() => {
                writeClock(jsonData);
            }, 1000)
        })
        .catch(err => {
            console.log(err)
        })
}

loadClock();

/* SVG daylight map */
function daylightMap() {
    writeWidget('daylight-map-interval', timeIntervalStr(daylightMapInterval))
    var DaylightMap, updateDateTime;

    DaylightMap = (function() {
        function DaylightMap(svg, date, options) {
            if (options == null) {
                options = {};
            }
            if (!((typeof SunCalc !== "undefined" && SunCalc !== null) && (typeof $ !== "undefined" && $ !== null) && (typeof d3 !== "undefined" && d3 !== null))) {
                throw new Error("Unmet dependency (requires d3.js, jQuery, SunCalc)");
            }
            if (!svg) {
                throw new TypeError("DaylightMap must be instantiated with a valid SVG");
            }
            this.options = {
                tickDur: options.tickDur || 400,
                shadowOpacity: options.shadowOpacity || 0.16,
                bgColorLeft: options.bgColorLeft || '#42448A',
                bgColorRight: options.bgColorRight || '#376281',
                lightsColor: options.lightsColor || '#FFBEA0',
                lightsOpacity: options.lightsOpacity || 0.5,
                sunOpacity: options.sunOpacity || 0.11
            };
            this.PRECISION_LAT = 1;
            this.PRECISION_LNG = 10;
            this.MAP_WIDTH = options.width || 1100;
            this.MAP_HEIGHT = this.MAP_WIDTH / 2;
            this.SCALAR_X = this.MAP_WIDTH / 360;
            this.SCALAR_Y = this.MAP_HEIGHT / 180;
            this.PROJECTION_SCALE = this.MAP_WIDTH / 6.25;
            this.WORLD_PATHS_URL = '/static/js/third-party/world-110m.json';
            this.CITIES_DATA_URL = '/static/js/third-party/cities-200000.json';
            this.svg = svg;
            this.isAnimating = false;
            this.cities = [];
            this.animInterval = null;
            this.currDate = date || new Date();
        }

        DaylightMap.prototype.colorLuminance = function(hex, lum) {
            var c, i, rgb;
            if (lum == null) {
                lum = 0;
            }
            c = null;
            i = 0;
            rgb = '#';
            hex = String(hex).replace(/[^0-9a-f]/gi, '');
            if (hex.length < 6) {
                hex = hex[0] + hex[0] + hex[1] + hex[1] + hex[2] + hex[2];
            }
            while (i < 3) {
                c = parseInt(hex.substr(i * 2, 2), 16);
                c = Math.round(Math.min(Math.max(0, c + c * lum), 255)).toString(16);
                rgb += ('00' + c).substr(c.length);
                i++;
            }
            return rgb;
        };

        DaylightMap.prototype.isDaylight = function(obj) {
            return obj.altitude > 0;
        };

        DaylightMap.prototype.isNorthSun = function() {
            return this.isDaylight(SunCalc.getPosition(this.currDate, 90, 0));
        };

        DaylightMap.prototype.getSunriseSunsetLatitude = function(lng, northSun) {
            var delta, endLat, lat, startLat;
            if (northSun) {
                startLat = -90;
                endLat = 90;
                delta = this.PRECISION_LAT;
            } else {
                startLat = 90;
                endLat = -90;
                delta = -this.PRECISION_LAT;
            }
            lat = startLat;
            while (lat !== endLat) {
                if (this.isDaylight(SunCalc.getPosition(this.currDate, lat, lng))) {
                    return lat;
                }
                lat += delta;
            }
            return lat;
        };

        DaylightMap.prototype.getAllSunPositionsAtLng = function(lng) {
            var alt, lat, peak, result;
            lat = -90;
            peak = 0;
            result = [];
            while (lat < 90) {
                alt = SunCalc.getPosition(this.currDate, lat, lng).altitude;
                if (alt > peak) {
                    peak = alt;
                    result = [peak, lat];
                }
                lat += this.PRECISION_LNG;
            }
            return result;
        };

        DaylightMap.prototype.getSunPosition = function() {
            var alt, coords, lng, peak, result;
            lng = -180;
            coords = [];
            peak = 0;
            while (lng < 180) {
                alt = this.getAllSunPositionsAtLng(lng);
                if (alt[0] > peak) {
                    peak = alt[0];
                    result = [alt[1], lng];
                }
                lng += this.PRECISION_LAT;
            }
            return this.coordToXY(result);
        };

        DaylightMap.prototype.getAllSunriseSunsetCoords = function(northSun) {
            var coords, lng;
            lng = -180;
            coords = [];
            while (lng < 180) {
                coords.push([this.getSunriseSunsetLatitude(lng, northSun), lng]);
                lng += this.PRECISION_LNG;
            }
            coords.push([this.getSunriseSunsetLatitude(180, northSun), 180]);
            return coords;
        };

        DaylightMap.prototype.lineFunction = d3.svg.line().x(function(d) {
            return d.x;
        }).y(function(d) {
            return d.y;
        }).interpolate('basis');

        DaylightMap.prototype.coordToXY = function(coord) {
            var x, y;
            x = (coord[1] + 180) * this.SCALAR_X;
            y = this.MAP_HEIGHT - (coord[0] + 90) * this.SCALAR_Y;
            return {
                x: x,
                y: y
            };
        };

        DaylightMap.prototype.getCityOpacity = function(coord) {
            if (SunCalc.getPosition(this.currDate, coord[0], coord[1]).altitude > 0) {
                return 0;
            }
            return 1;
        };

        DaylightMap.prototype.getCityRadius = function(population) {
            if (population < 200000) {
                return 0.3;
            } else if (population < 500000) {
                return 0.4;
            } else if (population < 100000) {
                return 0.5;
            } else if (population < 2000000) {
                return 0.6;
            } else if (population < 4000000) {
                return 0.8;
            } else {
                return 1;
            }
        };

        DaylightMap.prototype.getPath = function(northSun) {
            var coords, path;
            path = [];
            coords = this.getAllSunriseSunsetCoords(northSun);
            coords.forEach((function(_this) {
                return function(val) {
                    return path.push(_this.coordToXY(val));
                };
            })(this));
            return path;
        };

        DaylightMap.prototype.getPathString = function(northSun) {
            var path, pathStr, yStart;
            if (!northSun) {
                yStart = 0;
            } else {
                yStart = this.MAP_HEIGHT;
            }
            pathStr = "M 0 " + yStart;
            path = this.getPath(northSun);
            pathStr += this.lineFunction(path);
            pathStr += " L " + this.MAP_WIDTH + ", " + yStart + " ";
            pathStr += " L 0, " + yStart + " ";
            return pathStr;
        };

        DaylightMap.prototype.createDefs = function() {
            d3.select(this.svg).append('defs').append('linearGradient').attr('id', 'gradient').attr('x1', '0%').attr('y1', '0%').attr('x2', '100%').attr('y2', '0%');
            d3.select('#gradient').append('stop').attr('offset', '0%').attr('stop-color', this.options.bgColorLeft);
            d3.select('#gradient').append('stop').attr('offset', '100%').attr('stop-color', this.options.bgColorRight);
            d3.select(this.svg).select('defs').append('linearGradient').attr('id', 'landGradient').attr('x1', '0%').attr('y1', '0%').attr('x2', '100%').attr('y2', '0%');
            d3.select('#landGradient').append('stop').attr('offset', '0%').attr('stop-color', this.colorLuminance(this.options.bgColorLeft, -0.2));
            d3.select('#landGradient').append('stop').attr('offset', '100%').attr('stop-color', this.colorLuminance(this.options.bgColorRight, -0.2));
            d3.select(this.svg).select('defs').append('radialGradient').attr('id', 'radialGradient');
            d3.select('#radialGradient').append('stop').attr('offset', '0%').attr('stop-opacity', this.options.sunOpacity).attr('stop-color', "rgb(255, 255, 255)");
            return d3.select('#radialGradient').append('stop').attr('offset', '100%').attr('stop-opacity', 0).attr('stop-color', 'rgb(255, 255, 255)');
        };

        DaylightMap.prototype.drawSVG = function() {
            return d3.select(this.svg).attr('width', this.MAP_WIDTH).attr('height', this.MAP_HEIGHT).attr('viewBox', "0 0 " + this.MAP_WIDTH + " " + this.MAP_HEIGHT).append('rect').attr('width', this.MAP_WIDTH).attr('height', this.MAP_HEIGHT).attr('fill', "url(#gradient)");
        };

        DaylightMap.prototype.drawSun = function() {
            var xy;
            xy = this.getSunPosition();
            return d3.select(this.svg).append('circle').attr('cx', xy.x).attr('cy', xy.y).attr('id', 'sun').attr('r', 150).attr('opacity', 1).attr('fill', 'url(#radialGradient)');
        };

        DaylightMap.prototype.drawPath = function() {
            var path;
            path = this.getPathString(this.isNorthSun());
            return d3.select(this.svg).append('path').attr('id', 'nightPath').attr('fill', "rgb(0,0,0)").attr('fill-opacity', this.options.shadowOpacity).attr('d', path);
        };

        DaylightMap.prototype.drawLand = function() {
            return $.get(this.WORLD_PATHS_URL, (function(_this) {
                return function(data) {
                    var projection, worldPath;
                    projection = d3.geo.equirectangular().scale(_this.PROJECTION_SCALE).translate([_this.MAP_WIDTH / 2, _this.MAP_HEIGHT / 2]).precision(0.1);
                    worldPath = d3.geo.path().projection(projection);
                    d3.select(_this.svg).append('path').attr('id', 'land').attr('fill', 'url(#landGradient)').datum(topojson.feature(data, data.objects.land)).attr('d', worldPath);
                    return _this.shuffleElements();
                };
            })(this));
        };

        DaylightMap.prototype.drawCities = function() {
            return $.get(this.CITIES_DATA_URL, (function(_this) {
                return function(data) {
                    return data.forEach(function(val, i) {
                        var coords, id, opacity, radius, xy;
                        coords = [parseFloat(val[2]), parseFloat(val[3])];
                        xy = _this.coordToXY(coords);
                        id = "city" + i;
                        opacity = _this.getCityOpacity(coords);
                        radius = _this.getCityRadius(val[0]);
                        d3.select(_this.svg).append('circle').attr('cx', xy.x).attr('cy', xy.y).attr('id', id).attr('r', radius).attr('opacity', opacity * _this.options.lightsOpacity).attr('fill', _this.options.lightsColor);
                        return _this.cities.push({
                            title: val[1],
                            country: val[5],
                            latlng: coords,
                            xy: xy,
                            population: parseInt(val[0]),
                            id: id,
                            opacity: opacity
                        });
                    });
                };
            })(this));
        };

        DaylightMap.prototype.searchCities = function(str) {
            var cities;
            cities = _.filter(this.cities, function(val) {
                return val.title.toLowerCase().indexOf(str) === 0;
            });
            cities = _.sortBy(cities, function(val) {
                return val.population;
            });
            return cities.reverse();
        };

        DaylightMap.prototype.redrawSun = function(animate) {
            var curX, xy;
            xy = this.getSunPosition();
            curX = parseInt(d3.select("#sun").attr('cx'));
            if (animate && ((Math.abs(xy.x - curX)) < (this.MAP_WIDTH * 0.8))) {
                return d3.select("#sun").transition().duration(this.options.tickDur).ease('linear').attr('cx', xy.x).attr('cy', xy.y);
            } else {
                return d3.select("#sun").attr('cx', xy.x).attr('cy', xy.y);
            }
        };

        DaylightMap.prototype.redrawCities = function() {
            var k;
            k = 0;
            return this.cities.forEach((function(_this) {
                return function(val, i) {
                    var opacity;
                    opacity = _this.getCityOpacity(val.latlng);
                    if (val.opacity !== opacity) {
                        _this.cities[i].opacity = opacity;
                        k++;
                        return d3.select("#" + val.id).transition().duration(_this.options.tickDur * 2).attr('opacity', _this.options.lightsOpacity * opacity);
                    }
                };
            })(this));
        };

        DaylightMap.prototype.redrawPath = function(animate) {
            var nightPath, path;
            path = this.getPathString(this.isNorthSun(this.currDate));
            nightPath = d3.select('#nightPath');
            if (animate) {
                return nightPath.transition().duration(this.options.tickDur).ease('linear').attr('d', path);
            } else {
                return nightPath.attr('d', path);
            }
        };

        DaylightMap.prototype.redrawAll = function(increment, animate) {
            if (increment == null) {
                increment = 15;
            }
            if (animate == null) {
                animate = true;
            }
            this.currDate.setMinutes(this.currDate.getMinutes() + increment);
            this.redrawPath(animate);
            this.redrawSun(animate);
            return this.redrawCities();
        };

        DaylightMap.prototype.drawAll = function() {
            this.drawSVG();
            this.createDefs();
            this.drawLand();
            this.drawPath();
            this.drawSun();
            return this.drawCities();
        };

        DaylightMap.prototype.shuffleElements = function() {
            $('#land').insertBefore('#nightPath');
            return $('#sun').insertBefore('#land');
        };

        DaylightMap.prototype.animate = function(increment) {
            if (increment == null) {
                increment = 0;
            }
            if (!this.isAnimating) {
                this.isAnimating = true;
                return this.animInterval = setInterval((function(_this) {
                    return function() {
                        _this.redrawAll(increment);
                        return $(document).trigger('update-date-time', _this.currDate);
                    };
                })(this), this.options.tickDur);
            }
        };

        DaylightMap.prototype.stop = function() {
            this.isAnimating = false;
            return clearInterval(this.animInterval);
        };

        DaylightMap.prototype.init = function() {
            this.drawAll();
            return setInterval((function(_this) {
                return function() {
                    if (_this.isAnimating) {
                        return;
                    }
                    _this.redrawAll(1, false);
                    return $(document).trigger('update-date-time', _this.currDate);
                };
            })(this), 60000);
        };

        return DaylightMap;

    })();

    updateDateTime = function(date) {
        $('.curr-time').find('span').html(moment(date).format("HH:mm"));
        return $('.curr-date').find('span').text(moment(date).format("DD MMM"));
    };

    $(document).ready(function() {
        var map, svg;
        svg = document.getElementById('daylight-map');
        svg.innerHTML = "";
        map = new DaylightMap(svg, new Date());
        map.init();
        updateDateTime(map.currDate);
        $(document).on('update-date-time', function(date) {
            return updateDateTime(map.currDate);
        });
        $('.toggle-btn').on('click', function(e) {
            var $el;
            e.preventDefault();
            $el = $(this);
            return $el.toggleClass('active');
        });
        $('.js-skip').on('click', function(e) {
            var $el, animate;
            e.preventDefault();
            $el = $(this);
            animate = false;
            map.stop();
            $('.js-animate').removeClass('animating');
            if ($el.attr('data-animate')) {
                animate = true;
            }
            map.redrawAll(parseInt($(this).attr('data-skip')), animate);
            return updateDateTime(map.currDate);
        });
        return $('.js-animate').on('click', function(e) {
            var $el;
            $el = $(this);
            e.preventDefault();
            if ($el.hasClass('animating')) {
                $el.removeClass('animating');
                return map.stop();
            } else {
                $el.addClass('animating');
                return map.animate(10);
            }
        });
    });

}

daylightMap();
setInterval(() => {
    daylightMap();
}, daylightMapInterval)