window.addEventListener('load', function() {
    let ticking = false;

    const websiteZoomContainer = document.getElementById('website-zoom-container');

    function onScroll() {
        const scroll = window.scrollY || document.documentElement.scrollTop;
        const scrollMax = document.body.scrollHeight - window.innerHeight;
        const scaleStart = 0.3; // The min zoom size
        let scale;
        const scaleChangeRate = 0.7; // adjust as required.

        // calculate the end of scroll at which scaling should be complete (scale == 1)
        var zoomEnd = scaleChangeRate * scrollMax;
        if (scroll <= zoomEnd) {
            scale = (1 - scaleStart) * (scroll / zoomEnd) + scaleStart;
        } else {
            scale = 1;
        }

        if (scale === 1) {
            websiteZoomContainer.style.transform = 'none';
        } else {
            websiteZoomContainer.style.transform = 'scale(' + scale.toFixed(2) + ')';
        }
    }


    window.addEventListener('scroll', function() {
        if (!ticking) {
            window.requestAnimationFrame(function() {
                onScroll();
                ticking = false;
            });
        }
        ticking = true;
    });

    onScroll(); // Trigger scroll event to set initial values
});
