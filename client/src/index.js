var Calendar = tui.Calendar;
var calendarId = 'calendar';
var currentYear = new Date().getFullYear();
var currentMonth = new Date().getMonth();

var templates = {
  time: function (schedule) {
    return schedule.title + ' (' + ('<strong>' + moment(schedule.start.getTime()).format('HH:mm') + '</strong>' + ' ~ ' + '<strong>' + moment(schedule.end.getTime()).format('HH:mm') + '</strong>') + ')';
  }
};

var MONTHLY_CUSTOM_THEME = {
  // month header 'dayname'
  'month.dayname.height': '42px',
  'month.dayname.borderLeft': 'none',
  'month.dayname.paddingLeft': '8px',
  'month.dayname.paddingRight': '0',
  'month.dayname.fontSize': '13px',
  'month.dayname.backgroundColor': 'inherit',
  'month.dayname.fontWeight': 'normal',
  'month.dayname.textAlign': 'left',

  // month day grid cell 'day'
  'month.holidayExceptThisMonth.color': '#f3acac',
  'month.dayExceptThisMonth.color': '#bbb',
  'month.weekend.backgroundColor': '#fafafa',
  'month.day.fontSize': '16px',
};

function rfc3339(d) {
  function pad(n) {
    return n < 10 ? "0" + n : n;
  }

  function timezoneOffset(offset) {
    var sign;
    if (offset === 0) {
      return "Z";
    }
    sign = (offset > 0) ? "-" : "+";
    offset = Math.abs(offset);
    return sign + pad(Math.floor(offset / 60)) + ":" + pad(offset % 60);
  }

  return d.getFullYear() + "-" +
    pad(d.getMonth() + 1) + "-" +
    pad(d.getDate()) + "T" +
    pad(d.getHours()) + ":" +
    pad(d.getMinutes()) + ":" +
    pad(d.getSeconds()) +
    timezoneOffset(d.getTimezoneOffset());
}

function getSchedulePromise(year, month) {
  var apiUrl = "https://api-schedule.zychspace.com/schedule";
  var startDate = new Date(year, month, 1);
  startDate.setDate(startDate.getDate() - 7);
  var startDateStr = "startTime=" + encodeURIComponent(rfc3339(startDate));
  // Date change year and month automatically if month is 13
  var endDateStr = "endTime=" + encodeURIComponent(rfc3339(new Date(year, month + 1, 14)));

  var url = apiUrl + "?" + startDateStr + "&" + endDateStr;

  return axios.get(url);
}

function getCurrentMonth() {
  const d = new Date();
  const month = d.getMonth()+1;
  return month < 10 ? "0"+month : month;
}

function getCurrentYear() {
  const d = new Date();
  return d.getFullYear().toString();
}

function setCalendar(year, month) {
  document.getElementById('title').innerHTML = getCurrentYear() + "." + getCurrentMonth();

  document.getElementById('calendar').innerHTML = '';
  var calendar = new Calendar('#' + calendarId, {
    defaultView: 'month',
    scheduleView: ['time'],
    isReadOnly: true,
    useDetailPopup: true,
    template: templates,
    theme: MONTHLY_CUSTOM_THEME
  });

  // TODO: loading icon
  scheduleListPromise = getSchedulePromise(year, month);
  scheduleListPromise.then(resp => {
    calendarItemList = resp.data.map((data, idx) => ({
      id: idx + "",
      calendarId: calendarId,
      title: data.summary,
      category: 'time',
      start: data.start.dateTime,
      end: data.end.dateTime
    }));
    calendar.createSchedules(calendarItemList);
  });
}

window.onload = function() {
  setCalendar(currentYear, currentMonth);
}
