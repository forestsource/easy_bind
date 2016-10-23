$(function() {
  $("#toggle_zone").on("click", function() {
    $(".tab-item.active").removeClass("active")
    $("section.visible").removeClass("visible")
    $("#toggle_zone").addClass("active")
    $("#zone_area").addClass("visible")
  });
  $("#toggle_acl").on("click", function() {
    $(".tab-item.active").removeClass("active")
    $("section.visible").removeClass("visible")
    $("#toggle_acl").addClass("active")
    $("#acl_area").addClass("visible")
  });
  $("#toggle_key").on("click", function() {
    $(".tab-item.active").removeClass("active")
    $("section.visible").removeClass("visible")
    $("#toggle_key").addClass("active")
    $("#key_area").addClass("visible")
  });
  $("#toggle_option").on("click", function() {
    $(".tab-item.active").removeClass("active")
    $("section.visible").removeClass("visible")
    $("#toggle_option").addClass("active")
    $("#option_area").addClass("visible")
  });
});
