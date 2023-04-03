require("expose-loader?exposes=$,jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");

$(function(){
    $(".alert-dismissible").delay(2000).slideUp(200, function() {
        $(this).alert('close');
    });
})