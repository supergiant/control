// This file contains all the mess for rendering list views

$(function() {

  var table = $("table#item_list"),
      uiBasePath = table.data("ui-base-path"),
      apiListPath = table.data('api-list-path'),
      fields = table.data('fields-json'),
      thead = $('<thead>'),
      tbody = $('<tbody>'),
      paginator = $('ul.pagination');


  // Table headers
  var headerTrHtml = '<tr>';
  headerTrHtml += "<th></th>"; // for checkbox
  headerTrHtml += '<th>ID</th>';

  $.each(fields, function(fx, field) {
    headerTrHtml += '<th>' + field.title + '</th>'
  });

  headerTrHtml += '<th>Status</th>';
  headerTrHtml += '</tr>';

  thead.append($(headerTrHtml))
  table.append(thead);
  table.append(tbody);





  var currentPage = 1;
  var pmatch = window.location.search.match(/p=(\d+)/);
  if (pmatch) {
    currentPage = parseInt(pmatch[1]);
  }

  var limit = 25,
      offset = (currentPage - 1) * limit;

  var filterBits = window.location.search.match(/filter\.[^=]+=[^&]+/g);

  var url = apiListPath + "?limit=" + limit + "&offset=" + offset;
  if (filterBits) {
    url += "&" + filterBits.join('&');
  }

  var searchBar = $('input#searchbar'),
      searchIcon = $('#searchicon'),
      filterDiv = $('div#filters');

  searchBar.focusin(function() {
    searchIcon.css("color", "black");
  });
  searchBar.focusout(function() {
    searchIcon.css("color", "#bbb");
  });



  // Render the doo-dads

  function renderFilterPill(key, val) {
    var filterPill = $('<span class="search-filter"><strong>' + key + ':</strong> ' + val + '<span class="glyphicon glyphicon-remove"></span></span>');
    filterDiv.append(filterPill);
  }

  function padSearchBar() {
    searchBar.css("padding-left", filterDiv.width() + 10)
  }

  if (filterBits) {
    searchBar.prop("placeholder", "");

    filterBits.forEach(function(filterBit) {
      var filterMatch = filterBit.match(/filter\.([^=]+)=([^&]+)/);
      renderFilterPill(filterMatch[1], filterMatch[2])
    });

    padSearchBar();
  }



  // Adding filters
  //--------------------------------------------------------------------
  searchBar.keyup(function(event) {

    // enter key
    if (event.keyCode != 13) {
      return false;
    }

    var input = searchBar.val();

    if (input.length) {

      var keyVal = input.split(':');

      if (keyVal.length == 1) {
        keyVal = ["name", keyVal[0]];
      } else if (keyVal.length != 2) {
        return false;
      }

      var joiner;
      if (window.location.search.length) {
        joiner = "&";
      } else {
        joiner = "?";
      }

      // Just for the satisfaction
      searchBar.val("");
      searchBar.prop("placeholder", "");
      renderFilterPill(keyVal[0], keyVal[1]);
      padSearchBar();

      window.location.search += joiner + "filter." + keyVal[0] + "=" + keyVal[1];

    }

  });


  // Removing filters
  //--------------------------------------------------------------------
  $('span.search-filter .glyphicon-remove').click(function() {
    // Just for the satisfaction, remove the pill
    var pill = $(this).parent();
    pill.remove();
    padSearchBar();
    // Remove filter param, reload page
    var keyVal = pill.text().split(': ');

    var newQ = window.location.search.replace(RegExp("[\?&]filter." + keyVal[0] + "=" + keyVal[1]), "");

    // if (newQ.startsWith("&")) {
    //   newQ = newQ.substring(1);
    // }

    window.location.search = newQ;
  });




  // func
  function renderPagination(list) {
    var totalPages = Math.ceil(list.total / list.limit);

    if (totalPages > 1) {

      var newPaginator = paginator.clone();

      var renderPageLink = function(p) {
        if (p <= totalPages) {
          newPaginator.append($('<li' + (currentPage == p ? ' class="active"' : '') + '><a href="' + uiBasePath + '?p=' + p + '">' + p + '</a></li>'));
        }
      };
      var renderSeparator = function() {
        newPaginator.append($('<li class="disabled"><a href="#">...</a></li>'));
      };

      if (totalPages > 1) {
        if (currentPage > 1) {
          newPaginator.append($('<li><a href="' + uiBasePath + '?p=' + (currentPage - 1) + '" aria-label="Previous"><span aria-hidden="true">&laquo;</span></a></li>'));
        } else {
          newPaginator.append($('<li class="disabled"><a href="#" aria-label="Previous"><span aria-hidden="true">&laquo;</span></a></li>'));
        }
      }


      if (totalPages <= 6) {

        // Render all pages
        for (var p = 1; p <= totalPages; p++) {
          renderPageLink(p);
        }

      } else if (currentPage <= 3) { // and >= 7 total pages

        for (var p = 1; p <= (currentPage + 1); p++) {
          renderPageLink(p);
        }
        renderSeparator();
        renderPageLink(totalPages);

      } else if ((totalPages - currentPage) <= 3) { // and >= 7 total pages

        renderPageLink(1);
        renderSeparator();
        for (var p = currentPage - 1; p <= totalPages; p++) {
          renderPageLink(p);
        }

      } else {

        renderPageLink(1);
        renderSeparator();
        renderPageLink(currentPage - 1);
        renderPageLink(currentPage);
        renderPageLink(currentPage + 1);
        renderSeparator();
        renderPageLink(totalPages);

      }


      if (totalPages > 1) {
        if (currentPage < totalPages) {
          newPaginator.append($('<li><a href="' + uiBasePath + '?p=' + (currentPage + 1) + '" aria-label="Next"><span aria-hidden="true">&raquo;</span></a></li>'));
        } else {
          newPaginator.append($('<li class="disabled"><a href="#" aria-label="Next"><span aria-hidden="true">&raquo;</span></a></li>'));
        }
      }


      paginator.replaceWith(newPaginator);

    } else {

      paginator.empty();

    }
  }


  // func
  function renderItems(items) {
    // Remove any non-updated items (assuming they're deleted)
    $('tr.item[data-old-item="true"]').remove();
    // Mark existing items as old
    $("tr.item").attr("data-old-item", "true");



    $.each(items, function(i, item) {

      var previousTr = $('tr#item_' + item.id);

      // Gather field vals
      var fieldVals = [];
      $.each(fields, function(fx, field) {

        var val;

        if (field.type == "field_value") {
          var keys = field.field.split(".");
          val = item;
          $.each(keys, function(kx, key) {
            val = val[key];
          })
        } else if (field.type == "percentage") {

          var numerator = item[field.numerator_field],
              denominator = item[field.denominator_field],
              percent = Math.round(numerator * 100 / denominator);

          val = '<div class="progress" style="height: 10px; margin: 4px 0 0 0">' +
                '<div class="progress-bar progress-bar-info" role="progressbar" aria-valuenow="' + percent + '" aria-valuemin="0" aria-valuemax="100" style="width: ' + percent + '%">' +
                '</div>' +
                '</div>';
        }

        fieldVals.push(val);

      });

      // Build <tr> HTML
      var trHtml = '<tr class="item" data-old-item="false" id="item_' + item.id + '">';

      // Checkbox column (for batch actions)
      // var disabledCheckbox = item.status ? " disabled" : "";
      var disabledCheckbox = ""; // TODO we do need some type of disabling, but doing this prevents deleting hanging deploy
      trHtml += '<td class="item_selector"><input type="checkbox" id="' + item.id + '"' + disabledCheckbox + '></td>';

      // ID column
      trHtml += '<td><a href="' + uiBasePath + '/' + item.id + '">' + item.id + '</a></td>';

      // Attribute columns
      $.each(fieldVals, function(vx, val) {
        trHtml += '<td>' + val + '</td>';
      });

      // Status column
      var statusTd = '<td';
      if (item.status) {
        var color = item.status.description == "deleting" ? "danger" : "info";
        statusTd += ' class="text-' + color + '">';

        statusTd += '<span>' + item.status.description + '</span>';

        // Loader
        statusTd += '<div id="circleG"><div id="circleG_1" class="circleG"></div><div id="circleG_2" class="circleG"></div><div id="circleG_3" class="circleG"></div></div>'

        statusTd += '</td>'

      } else if (item.passive_status) {

        var color = item.passive_status_okay ? "success": "warning";
        statusTd += ' style="opacity: 0.7" class="text-' + color + '">' + item.passive_status + '</td>';

      } else {
        statusTd += '></td>';
      }
      trHtml += statusTd;

      trHtml += '</tr>';


      var newTr = $(trHtml);

      if (previousTr.length) {

        // Don't replace (we don't want flashing) if it hasn't changed.
        // However, do update this attribute so it's not deleted.
        if (previousTr[0].innerHTML == newTr[0].innerHTML) {
          previousTr.attr("data-old-item", "false")
        } else {
          // console.log("REPLACING");
          previousTr.replaceWith(newTr);
        }
      } else {
        // console.log("APPENDING")
        tbody.append($(trHtml));
      }

    });
  }



  // func
  function renderList() {
    $.ajax({
      url: url,
      beforeSend: function(xhr){
        xhr.setRequestHeader('Authorization', 'SGAPI session="' + getCookie('supergiant_session') + '"');
      },
      success: function(data) {
        var list = JSON.parse(data);
        var items = list.items;
        renderItems(items);
        renderPagination(list);

        // see resource_graphs.js for example on how to use
        if (window.postListRender) {
          window.postListRender(list);
        }

      },
      error: function(data) {
        console.warn(data)

        // Session no longer exists, force login
        if (data.status == 401) {
          location.reload()
        }
      }
    });
  }




  renderList();
  setInterval(renderList, 2000);



  var actionLinks = $('a[data-action-path]'),
      batchActionLinks = $('a[data-batch-action-path]'),
      confirmModal = $('.modal#confirm_action'),
      modalBody = $('.model-body'),
      modalActionName = $("strong#modal_action_name"),
      modalList = $("ul#modal_list"),
      modalConfirmBtn = $('#confirm_action_btn'),
      selectedItemIDs = [];



  tbody.on("click", "td.item_selector", function(event) {
    var input = $(this).children('input'),
        itemID = input[0].id;

    if (event.target.tagName == "INPUT") {
      // do its natural thing...
    } else {
      input.prop('checked', !input.prop('checked')) // toggle
    }

    selectedItemIDs = [];
    $("td.item_selector > input:checked").each(function(chx, checkbox) {
      selectedItemIDs.push(checkbox.id);
    });

    // Single actions
    if (selectedItemIDs.length == 1) {
      actionLinks.each(function(lx, actionLink) {
        var link = $(actionLink),
            path = uiBasePath + "/" + itemID + link.data("action-path");
        link.attr("href", path);
      });
      actionLinks.removeAttr('disabled');
      actionLinks.removeClass('disabled');
    } else {
      actionLinks.attr("href", "#");
      actionLinks.attr('disabled', 'disabled');
      actionLinks.addClass('disabled');
    }

    // Batch actions
    if (selectedItemIDs.length > 0) {
      batchActionLinks.removeAttr('disabled');
      batchActionLinks.removeClass('disabled');
    } else {
      batchActionLinks.attr('disabled', 'disabled');
      batchActionLinks.addClass('disabled');
    }
  });


  // Display and render items in modal on action link click
  batchActionLinks.on('click', function() {
    var link = $(this);

    if (link.hasClass("disabled")) {
      return false;
    }

    var actionName = link.text();
    modalActionName.text(actionName);


    // TODO
    if (actionName == "Delete") {
      modalConfirmBtn.addClass("btn-danger");
      modalActionName.addClass("text-danger");
    } else {
      modalConfirmBtn.removeClass("btn-danger");
      modalActionName.removeClass("text-danger");
    }


    modalList.empty();
    $.each(selectedItemIDs, function(selx, id) {
      modalList.append($("<li>" + id + "</li>"));
    });

    modalConfirmBtn.text(actionName);
    modalConfirmBtn.data("batch-action-path", link.data("batch-action-path"));

    confirmModal.modal();

    return false;
  });


  // On action confirmation, iterate through items and make requests
  modalConfirmBtn.on('click', function() {
    // TODO
    var alerted = false;

    $.each(selectedItemIDs, function(selx, id) {

      console.log(uiBasePath + "/" + id + modalConfirmBtn.data("batch-action-path"));

      $.ajax({
        type: "PUT",
        beforeSend: function(xhr){
          xhr.setRequestHeader('Authorization', 'SGAPI session="' + getCookie('supergiant_session') + '"');
        },
        url: uiBasePath + "/" + id + modalConfirmBtn.data("batch-action-path"),
        error: function(data) {
          if (!alerted) {
            alert(data.responseText);
          }
          alerted = true;
        },
        success: function() {

          if ((selx + 1) == selectedItemIDs.length) {
            renderList();
          }

        }
      });
    });

    confirmModal.modal("hide");
  });


  // 1. on checkbox click
  //
  // 2. gather selected item IDs, and store globally
  //
  // 3. if empty, un-disable action links in drop down
  //
  // ----
  //
  // 1. on action link click
  //
  // 2. modal pops up, says action name with list of all selected items
  //
  // ----
  //
  // 1. on confirmation, iterate through all items and make action request
  //
  // 2. wait for response (since async) for every item
  //
  // 3. renderList()
  //
  // 4. close modal




  var logDiv = $('#log');
  if (logDiv.length) {

    var renderLogs = function() {
      $.get({
        beforeSend: function(xhr){
          xhr.setRequestHeader('Authorization', 'SGAPI session="' + getCookie('supergiant_session') + '"');
        },
        url: "/api/v0/log",
        success: function(data) {

          data = data.replace(/[\x00-\x7F]\[\d+mINFO[\x00-\x7F]\[0m/g, "<span class='text-info'>INFO</span> ")
          data = data.replace(/[\x00-\x7F]\[\d+mWARN[\x00-\x7F]\[0m/g, "<span class='text-warning'>WARN</span> ")
          data = data.replace(/[\x00-\x7F]\[\d+mERRO[\x00-\x7F]\[0m/g, "<span class='text-danger'>ERRO</span> ")
          data = data.replace(/[\x00-\x7F]\[\d+mDEBU[\x00-\x7F]\[0m/g, "<span class='text-muted'>DEBU</span> ")

          logDiv.html(data);

        }
      });
    };

    renderLogs();

    // Scroll to bottom (like a terminal)
    logDiv.animate({
      scrollTop: logDiv.height()
    }, 200);

    setInterval(renderLogs, 2000);
  }


});
