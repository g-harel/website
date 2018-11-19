// Default hash will be active if the current one is empty/invalid.
var defaultHash = "#projects";

// This class is added/removed from dom elements to control styling.
var activeClass = "hash-active";

// Page contents are only made visible after tabbed navigation has
// loaded to prevent flash of empty content.
document.body.style.opacity = 0;
document.addEventListener("DOMContentLoaded", () => {
    document.body.style.opacity = 1;
});

function isValidHash(hash) {
    if (!hash) {
        return false;
    }
    if ([defaultHash, "#contributions"].indexOf(hash) === -1) {
        return false;
    }
    return true;
}

function updateActive() {
    var hash = document.location.hash;

    if (!isValidHash(hash)) {
        document.location.hash = defaultHash;
        // Because this function is listening to hash changes,
        // the first invocation must not run if it reassigns
        // the hash value.
        return;
    }

    var formattedHash = hash.replace(/^#/, "");
    if (formattedHash === "") {
        return;
    }

    // Previously active elements have the active class removed.
    var oldActive = document.querySelectorAll("." + activeClass);
    for (var i = 0; i < oldActive.length; i++) {
        oldActive[i].classList.remove(activeClass);
    }

    // Newly active elements have the active class added.
    var newActive = document.querySelectorAll("." + formattedHash);
    for (var i = 0; i < newActive.length; i++) {
        newActive[i].classList.add(activeClass);
    }
}

window.addEventListener("hashchange", updateActive);
document.addEventListener("DOMContentLoaded", updateActive);

if (!isValidHash(document.location.hash)) {
    document.location.hash = defaultHash;
}
