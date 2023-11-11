window.addEventListener('load', function() {
    var websiteZoomContainer = document.getElementById('website-zoom-container');

    function onScroll() {
        var scroll = window.scrollY || document.documentElement.scrollTop;
        var scrollMax = document.body.scrollHeight - window.innerHeight;
        var scaleStart = 0.3; // The min zoom size
        var scale = scaleStart;
        var bottom = 0;
        var scaleChangeRate = 0.7; // adjust as required.

        // calculate the end of scroll at which scaling should be complete (scale == 1)
        var zoomEnd = scaleChangeRate * scrollMax;
        if (scroll <= zoomEnd) {
            scale = (1 - scaleStart) * (scroll / zoomEnd) + scaleStart;
        } else {
            scale = 1;
            bottom = ((scroll - zoomEnd) / (scrollMax - zoomEnd)) * 100;
        }

        if (scale === 1) {
            websiteZoomContainer.style.transform = 'none';
        } else {
            websiteZoomContainer.style.transform = 'scale(' + scale.toFixed(2) + ')';
        }
    }

    window.addEventListener('scroll', onScroll);

    onScroll(); // Trigger scroll event to set initial values
});
