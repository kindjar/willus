var WILLUS = (function (my) {

  my.minutelyWidth = 200;
  my.minutelyHeight = 50;
  my.hourlyWidth = 400;
  my.hourlyHeight = 50;
  my.dailyPrecipitationWidth = 760;
  my.dailyPrecipitationHeight = 100;
  my.dailyTemperatureWidth = 800;
  my.dailyTemperatureHeight = 300;

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

  my.makePrecipChart = function (selector, data, width, height) {
    var barWidth = width / data.length;
    var y = d3.scale.linear()
        .range([0, height / d3.max(data, function(d) { return d[0]; })]);

    d3.select(selector + ' .y-axis').style('height', height + "px");
    d3.select(selector + ' .y-axis .top').
        text(my.precipIntensityLabel(
            d3.max(data, function(d) { return d[0] })));
    d3.select(selector + ' .y-axis .bottom').
        text(my.precipIntensityLabel(
            d3.min(data, function(d) { return d[0] })));

    return d3.select(selector + ' .area').
      style('width', width + "px").
      style('height', height + "px").
      selectAll("div").
        data(data).
      enter().append("div").
        classed('bar', true).
        style("height", function(d) {
          return y(d[0]) + "px";
        }).
        style("width", barWidth + "px").
        style('left', function(d, i) { return ((i * barWidth) + "px"); }).
        style('opacity', function(d) { return d[1]; }).
        attr('title', function(d) { return "" + d[0] + " (" + d[1] + ")"; });
  };

  my.makeTempChart = function (selector, data, width, height) {
    var margin = {top: 0, right: 0, bottom: 30, left: 50},
        interiorWidth = width - margin.left - margin.right,
        interiorHeight = height - margin.top - margin.bottom;
    var x = d3.time.scale().range([0, interiorWidth]);
    var y = d3.scale.linear().range([interiorHeight, 0]);
    var xAxis = d3.svg.axis().scale(x).
        orient("bottom").ticks(data.length / 2).
        tickFormat(d3.time.format("%-m/%d"));
    var yAxis = d3.svg.axis().scale(y).
        orient("left").ticks(3);
    var valueline = d3.svg.line().
        interpolate("basis").
        x(function(d) { return x(d[0]); }).
        y(function(d) { return y(d[1]); });

    // Scale the range of the data
    x.domain(d3.extent(data, function(d) { return d[0]; }));
    y.domain([
      d3.min(data, function(d) { return d[1]; }),
      d3.max(data, function(d) { return d[1]; })
    ]);

    var svg = d3.select(selector + ' .area')
        .append("svg")
            .attr("width", width)
            .attr("height", height)
        .append("g")
            .attr("transform",
                  "translate(" + margin.left + "," + margin.top + ")");

    // Add the valueline path.
    svg.append("path")
        .attr("class", "line")
        .attr("d", valueline(data));

    // Add the X Axis
    svg.append("g")
        .attr("class", "x axis")
        .attr("transform", "translate(0," + interiorHeight + ")")
        .call(xAxis);

    // Add the Y Axis
    svg.append("g")
        .attr("class", "y axis")
        .call(yAxis);
    return svg;
  }

  return my;
}(WILLUS || {}));
