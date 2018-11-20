function stopPropagationForScrollable() {
    var scrollable = document.querySelectorAll(".content");

    function captureScroll(e) {
        e.stopPropagation();
    }

    for (var i = 0; i < scrollable.length; i ++) {
        scrollable[i].addEventListener("wheel", captureScroll, {passive: true});
    }
}

function pulseShadeWhenScroll() {
    var shadeDelay = 200;
    var shade = document.querySelector(".shade");
    var timer = null;

    function pulseShade() {
        shade.style.opacity = 1;
        clearTimeout(timer);
        timer = setTimeout(() => {
            shade.style.opacity = 0;
        }, shadeDelay)
    }

    document.body.addEventListener("wheel", pulseShade, {passive: true});
}

document.addEventListener('DOMContentLoaded', function() {
    stopPropagationForScrollable();
    pulseShadeWhenScroll();
});
