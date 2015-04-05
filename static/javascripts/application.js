var WILLUS = (function (my) {

  my.minutelyWidth = 200;
  my.minutelyHeight = 50;
  my.hourlyWidth = 400;
  my.hourlyHeight = 50;
  my.dailyWidth = 400;
  my.dailyHeight = 100;

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

  return my;
}(WILLUS || {}));
