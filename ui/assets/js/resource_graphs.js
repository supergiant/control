var graphContainer = $('#two_graphs'),
    cpuGraph = d3.select('svg#cpu_graph'),
    ramGraph = d3.select('svg#ram_graph'),
    podLabelDiv = $('#pod_labels'),
    labelDivContainer = podLabelDiv.parent();



var marginTop = 14,
    marginBottom = 18,
    marginLeft = 0,
    marginRight = 42;



var podColors = {};

function getColor(podName) {
  var color = podColors[podName];
  if (!color) {
    color = podColors[podName] = {
      hue: Math.round(Math.random() * 360),
      lightness: 60 // 50 + Math.round(Math.random() * 20) // 50 - 70
    }
  }
  return color
}






function renderGraph(name, graph, podMetrics, yTickFormatter) {

  // Clear
  graph.selectAll("*").remove();

  podLabelDiv.empty();






  var height = graph.node().height.baseVal.value - (marginTop + marginBottom);
      // width = graph.node().width.baseVal.value - (marginLeft + marginRight);


  // NOTE
  //
  //
  //
  // For whatever reason we cannot get the correct width of the graphs on load.
  //
  //
  //
  //
  //
  var width = ($('#content').parent().width() / 2) - (marginLeft + marginRight);




  // set the ranges
  var x = d3.scaleTime().range([0, width]);
  var y = d3.scaleLinear().range([height, 0]);


  // Scale the range of the data
  var allValues = [];
  for (var podName in podMetrics) {
    var data = podMetrics[podName];
    // format timestamps
    data.forEach(function(d) {
      d.timestamp = d3.isoParse(d.timestamp);
    });
    allValues = allValues.concat(data);
  }
  x.domain(d3.extent(allValues, function(d) { return d.timestamp; }));
  y.domain([0, d3.max(allValues, function(d) { return d.value; }) * 1.2]); // just slightly larger than max





  // Title
  graph.append("text")
       .style("fill", "black")
       .style("font-size", "12px")
       .style("font-family", "'Open Sans', sans-serif")
       .attr("y", 10)
       .text(name);


  // Graph container
  var el = graph.append("g")
       .attr("transform", "translate(" + marginLeft + "," + marginTop + ")");




  for (var podName in podMetrics) {
    // Get data
    var color = getColor(podName);
    var data = podMetrics[podName];
    if (!data) {
      continue;
    }

    // area under the line
    var area = d3.area()
        .x(function(d) { return x(d.timestamp); })
        .y0(height)
        .y1(function(d) { return y(d.value); });
    // Add the line fill
    el.append("path")
      .datum(data)
      .style("fill", "hsla(" + color.hue + ",65%," + (color.lightness + 20) + "%,0.5)")
      .style("stroke", "none")
      .attr("d", area);

  }




  for (var podName in podMetrics) {
    // Get data
    var color = getColor(podName);
    var data = podMetrics[podName];
    if (!data) {
      continue;
    }

    var label = $('<span style="font-family: Open Sans, sans-serif; font-size: 11px; color: hsla(' + color.hue + ',65%,' + color.lightness + '%,1); padding-right: 8px">' + podName + '</span>')
    podLabelDiv.append(label);







    // define the line
    var valueline = d3.line()
        .x(function(d) { return x(d.timestamp); })
        .y(function(d) { return y(d.value); });
    // Add the valueline path.
    el.append("path")
        .data([data])
        .attr("class", "graph-line")
        .attr("d", valueline)
        .style("stroke", "hsla(" + color.hue + ",65%," + color.lightness + "%,0.9)")

  }


  // Add the X Axis
  el.append("g")
      .attr("class", "x axis")
      .attr("transform", "translate(0," + height + ")")
      .call(d3.axisBottom(x).ticks(6))


  // Add the Y Axis
  var yAxis = d3.axisRight(y)
               .ticks(3)
               .tickSize(width);
  if (yTickFormatter) {
    yAxis = yAxis.tickFormat(yTickFormatter);
  }

  var gy = el.append("g")
    .attr("class", "y axis")
    .call(yAxis);

  gy.selectAll("g").filter(function(d) { return d; })
    .classed("minor", true);

  gy.selectAll("text")
    .style("text-anchor", "end")
    .style("text-shadow", "2px 2px white")
    .attr("dx", -5)
    .attr("dy", -4);

}



// This is called by renderList() in index.js
window.postListRender = function(list) {


  var podCPUMetrics = {},
      podRAMMetrics = {};


  $.each(list.items, function(i, kubeResource) {
    if (kubeResource.kind == "Pod" && kubeResource.extra_data) {
      var fullName = kubeResource.namespace + "::" + kubeResource.name;
      podCPUMetrics[fullName] = kubeResource.extra_data.metrics.cpu_usage;
      podRAMMetrics[fullName] = kubeResource.extra_data.metrics.ram_usage;
    }
  });





  var maxDataPoints = 0;
  for (var podName in podCPUMetrics) {
    maxDataPoints = Math.max(maxDataPoints, podCPUMetrics[podName].length);
  }
  for (var podName in podRAMMetrics) {
    maxDataPoints = Math.max(maxDataPoints, podRAMMetrics[podName].length);
  }




  if (maxDataPoints < 2) {
    graphContainer.hide();
    labelDivContainer.hide();
    return
  }
  renderGraph("CPU (millicores)", cpuGraph, podCPUMetrics, function(millicores) {
    return d3.format(".1s")(millicores);
  });
  renderGraph("RAM (MiB)", ramGraph, podRAMMetrics, function(bytes) {
    return d3.format(".1s")(Math.ceil(bytes / 1048576));
  });
  graphContainer.show();
  labelDivContainer.show();



};
