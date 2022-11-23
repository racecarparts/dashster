const calendarInterval = 3600000 // every hour

function generateCalendar(d, id) {
  function monthDays(month, year) {
    var result = [];
    var days = new Date(year, month, 0).getDate();
    for (var i = 1; i <= days; i++) {
      result.push(i);
    }
    return result;
  }
  Date.prototype.monthDays = function() {
    var d = new Date(this.getFullYear(), this.getMonth() + 1, 0);
    return d.getDate();
  };
  var details = {
    // totalDays: monthDays(d.getMonth(), d.getFullYear()),
    totalDays: d.monthDays(),
    weekDays: ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa'],
    months: ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'],
  };

  var start = new Date(d.getFullYear(), d.getMonth()).getDay();
  var cal = [];

  var day = 1;
  var now = new Date();
  for (var i = 0; i <= 6; i++) {
    cal.push(['<tr>']);
    for (var j = 0; j < 7; j++) {
      if (i === 0) {
        cal[i].push('<td>' + details.weekDays[j] + '</td>');
      } else if (day > details.totalDays) {
        cal[i].push('<td>&nbsp;</td>');
      } else {
        if (i === 1 && j < start) {
          cal[i].push('<td>&nbsp;</td>');
        } else {
          var todaySpan = '<span class="'
          if (d.getMonth() === now.getMonth() && day === now.getDate()) {
            todaySpan = todaySpan + 'text-dark bg-light';
          }
          todaySpan = todaySpan + '">' + day++ + '</span>'
          cal[i].push('<td class="day">' + todaySpan + '</td>');
        }
      }
    }
    cal[i].push('</tr>');
  }

  // month header
  cal.unshift([
    '<tr><td class="text-center" colspan="7">' + 
    details.months[d.getMonth()] +
    '</td></tr>'
  ]);

  cal = cal.reduce(function(a, b) {
    return a.concat(b);
  }, []).join('');
  $('#' + id).append(cal);
}

function writeCalendar() {
    var currentDate = new Date();

    if (currentDate.getMonth() === 0) {
      currentDate = new Date(currentDate.getFullYear() - 1, 11);
      // generateCalendar(currentDate, 'prev');
    } else {
      currentDate = new Date(currentDate.getFullYear(), currentDate.getMonth() - 1)
      // generateCalendar(currentDate);
    }
    generateCalendar(currentDate, 'calendar-table-prev');

    currentDate = new Date();
    generateCalendar(currentDate, 'calendar-table-cur');

    if (currentDate.getMonth() === 11) {
      currentDate = new Date(currentDate.getFullYear() + 1, 0);
    } else {
      currentDate = new Date(currentDate.getFullYear(), currentDate.getMonth() + 1)
    }
    generateCalendar(currentDate, 'calendar-table-next');
}

writeCalendar()
setInterval(() => {
  writeCalendar()
}, calendarInterval)