$(function() {
  $("div.markdown").each(function(e) {
    $(this).html(marked($(this).text()));
  });

  $("form#dict-search").click(function(e) {
    e.preventDefault();
    var data = $(this).serialize();
    $.ajax({type: 'POST', url: '/dict', data: data}).done(function(rst) {
      if (rst) {
        $("div#dict-results").html(rst.map(function(it) {
          return it.data;
        }).join("<br/>"));
      } else {
        alert("没有结果")
      }
    });
  });

});
