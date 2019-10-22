$(document).ready(function () {
  $(".navbar-burger").click(function () {
    $(".button").toggleClass("is-active");
    $(".navbar-burger").toggleClass("is-active");
    $(".navbar-menu").toggleClass("is-active");
    $(".navbar-hidden").toggleClass("is-active");

    if ($(".is-sidebar-menu")[0].offsetWidth == 0 ||
      $(".is-sidebar-menu").hasClass("mobile-sidebar")) {
      $(".is-main-content").toggleClass("is-hidden-content");
      $(".is-sidebar-menu").toggleClass("is-hidden-mobile");
      $(".is-sidebar-menu").toggleClass("mobile-sidebar");
    }
  });
});
