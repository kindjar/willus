var WILLUS = (function (my) {

  my.minutelyWidth = 200;
  my.minutelyHeight = 150;
  my.hourlyWidth = 760;
  my.hourlyHeight = 150;
  my.dailyPrecipitationWidth = 760;
  my.dailyPrecipitationHeight = 150;
  my.dailyTemperatureWidth = 760;
  my.dailyTemperatureHeight = 300;
  my.minimumPrecipIntensityMax = 0.4;

  my.uniq = function(array) {
    return array.filter(function (value, index, array) {
      return array.indexOf(value) === index;
    });
  }

  my.parseDate = d3.time.format("%Y-%m-%dT%H:%M:%S%Z").parse;
  my.roundTimeTo = function(t, h, m, s) {
    t.setHours(h || 0);
    t.setMinutes(m || 0);
    t.setSeconds(s || 0);
    t.setMilliseconds(0);
    return t;
  }
  my.asTempMinTime = function(s) {
    return my.roundTimeTo(my.parseDate(s), 4);
  };
  my.asTempMaxTime = function(s) {
    return my.roundTimeTo(my.parseDate(s), 16);
  };

  my.precipIntensityLabel = function(inchesPerHour) {
    var label;
    if (inchesPerHour > 0.4 ) {
      label = "heavy";
    } else if (inchesPerHour > 0.1 ) {
      label = "moderate";
    } else if (inchesPerHour >= 0.017 ) {
      label = "light";
    } else if (inchesPerHour >= 0.002) {
      label = "very light";
    } else {
      label = "none";
    }
    return label;
  }

  my.makePrecipChart = function (selector, data, width, height, opts = {}) {
    var margin = {top: 10, right: 0, bottom: 25, left: 50},
        interiorWidth = width - margin.left - margin.right,
        interiorHeight = height - margin.top - margin.bottom,
        dates = my.uniq(data.map(
          function(d) { return d[0]; }
        ));
    var maxPrecip = Math.max(
      d3.max(data, function(d) { return d[1]; }), my.minimumPrecipIntensityMax
    );
    var timeFormatter = d3.time.format(opts['timeFormat'] || "%-m/%d");
    // var x = d3.time.scale().range([0, interiorWidth]);
    var x = d3.scale.ordinal().
      domain(data.map(function(d) { return d[0]; })).
      rangeRoundBands([0, interiorWidth], 0.1);
    var y = d3.scale.linear().range([interiorHeight, 0]).domain([0, maxPrecip]);
    var xAxis = d3.svg.axis().scale(x).
        orient("bottom").
        tickFormat(timeFormatter);

    if (opts['xLabelPeriod']) {
      var remainder = Math.floor(opts['xLabelPeriod'] / 2);
      xAxis.tickValues(dates.filter(function(v,i,a) {
        return i % opts['xLabelPeriod'] == remainder;
      }));
    } else {
      xAxis.tickValues(dates);
    }

    var yAxis = d3.svg.axis().scale(y).
        orient("left").
        // tickFormat(my.precipIntensityLabel).
        ticks(3);

    var chart = d3.select(selector + ' .area').
      append("svg").
        style("width", width).
        style("height", height).
      append("g").
        attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    chart.append("g").
      attr("class", "x axis").
      attr("transform", "translate(0," + interiorHeight + ")").
      call(xAxis);

    chart.append("g").
      attr("class", "y axis").
      call(yAxis);

    chart.selectAll(".bar").
      data(data).
      enter().append("rect").
        attr("class", "bar").
        attr("x", function(d, i) { return x(d[0]); }).
        attr("y", function(d) { return y(d[1]); }).
        attr("height", function(d) { return interiorHeight - y(d[1]); }).
        attr("width", x.rangeBand()).
        attr("opacity", function(d) { return d[2]; }).
        attr("title", function(d) { return "" + d[1] + '" (' + d[2] + "%)"; });
  };

  my.makeTempChart = function (selector, data, width, height) {
    var margin = {top: 0, right: 0, bottom: 30, left: 50},
        interiorWidth = width - margin.left - margin.right,
        interiorHeight = height - margin.top - margin.bottom,
        tempMarginDegrees = 2,
        dates = my.uniq(data.map(
          function(d) { return my.roundTimeTo(new Date(d[0].getTime()), 12); }
        ));
    var x = d3.time.scale().range([0, interiorWidth]);
    var y = d3.scale.linear().range([interiorHeight, 0]);
    var xAxis = d3.svg.axis().scale(x).
        orient("bottom").
        ticks(dates.length).
        tickValues(dates).
        tickFormat(d3.time.format("%-m/%d"));
    var yAxis = d3.svg.axis().scale(y).
        orient("left").
        tickFormat(function(d) { return "" + d + "\xB0"; });
    var valueline = d3.svg.line().
        interpolate("monotone").
        x(function(d) { return x(d[0]); }).
        y(function(d) { return y(d[1]); });

    // Scale the range of the data
    x.domain(d3.extent(data, function(d) { return d[0]; }));
    y.domain([
      d3.min(data, function(d) { return d[1]; }) - tempMarginDegrees,
      d3.max(data, function(d) { return d[1]; }) + tempMarginDegrees
    ]);

    var svg = d3.select(selector + ' .area').
      append("svg").
        attr("width", width).
        attr("height", height).
      append("g").
        attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    // Add the valueline path.
    svg.append("path").
        attr("class", "line").
        attr("d", valueline(data));

    // Add the X Axis
    svg.append("g").
        attr("class", "x axis").
        attr("transform", "translate(0," + interiorHeight + ")").
        call(xAxis);

    // Add the Y Axis
    svg.append("g").
        attr("class", "y axis").
        call(yAxis);
    return svg;
  }

  return my;
}(WILLUS || {}));
