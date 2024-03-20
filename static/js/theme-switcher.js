// leht laetakse sujuvalt sisse
document.addEventListener('DOMContentLoaded', () => {
    document.body.classList.add('loaded');

    const savedTheme = localStorage.getItem('theme') || 'light';
    document.body.setAttribute("data-theme", savedTheme);
});


// Need kaks järgmist koodi lõiku hoolitsevad lehekülje dark/light tekkimise ja püsimise eest

document.getElementById("themeToggle").addEventListener("click", function() {
    let currentTheme = document.body.getAttribute("data-theme") === "dark" ? "light" : "dark";
    document.body.setAttribute("data-theme", currentTheme);
    localStorage.setItem('theme', currentTheme); // Salvesta teema eelistus
});

// Lae teema eelistus, kui leht laaditakse
document.addEventListener('DOMContentLoaded', () => {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.body.setAttribute("data-theme", savedTheme);
});

// document.getElementById("themeToggle").addEventListener("click", function() {
//     if (document.body.getAttribute("data-theme") === "dark") {
//         document.body.setAttribute("data-theme", "light");
//     } else {
//         document.body.setAttribute("data-theme", "dark");
//     }
// });
