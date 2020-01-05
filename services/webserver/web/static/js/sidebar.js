$('.menu li').on('click', 'a', function (e) {
  if ($(this).parent().children('ul').length) {
    e.preventDefault();
    $(this).addClass('active');
    $(this).parent().children('ul').slideDown();
    setTimeout(function () {
      if ($.fn.matchHeight) {
        $.fn.matchHeight._update();
        $.fn.matchHeight._maintainScroll = true;
      }
    }, 500);
  }
});

/* Close the expanded menu */
$('.menu li').on('click', 'a.active', function (e) {
  if ($(this).parent().children('ul').length) {
    e.preventDefault();
    $(this).removeClass('active');
    $(this).parent().children('ul').slideUp();
    setTimeout(function () {
      if ($.fn.matchHeight) {
        $.fn.matchHeight._update();
        $.fn.matchHeight._maintainScroll = true;
      }
    }, 500);
  }
});
