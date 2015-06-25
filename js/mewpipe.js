var roles = {
    "148282": {
        name: "Arnaud-Pierre BOYER",
        title: "CEO"
    },
    "169174": {
        name: "Gaël DEMETTE",
        title: "CTO"
    },
    "170768": {
        name: "Alexandre VAAST",
        title: "Tech. sales"
    },
    "202429": {
        name: "Robin LEBERT",
        title: "Front. engineer"
    },
    "207910": {
        name: "Mohamed Amir BEN SLAMIA",
        title: "Infra. engineer"
    }
};
function SlideChangedHandler (event) {

    if (!event.currentSlide.dataset.author && event.currentSlide.parentNode.nodeName == 'SECTION' && event.currentSlide.parentNode.dataset.author) {
        event.currentSlide.dataset.author = event.currentSlide.parentNode.dataset.author;
    }

    if (event.currentSlide.dataset.author) {
        var author = roles[event.currentSlide.dataset.author]
        document.getElementById("author-card").removeAttribute("aria-disabled");
        document.getElementById("author-photo").setAttribute("src", "img/" + event.currentSlide.dataset.author + ".jpg");
        document.getElementById("author-name").innerHTML = author.name;
        document.getElementById("author-title").innerHTML = author.title;

    }
    else if (event.previousSlide && event.previousSlide.dataset.author) {
        document.getElementById("author-card").setAttribute("aria-disabled", "true");
    }

}


Reveal.addEventListener('slidechanged',  SlideChangedHandler);
Reveal.addEventListener( 'ready', SlideChangedHandler);