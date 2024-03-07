document.getElementById("themeToggle").addEventListener("click", function() {
    if (document.body.getAttribute("data-theme") === "dark") {
        document.body.setAttribute("data-theme", "light");
    } else {
        document.body.setAttribute("data-theme", "dark");
    }
});
