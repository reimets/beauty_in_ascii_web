// leht laetakse sujuvalt sisse ja laeb lehe eelistuse 
document.addEventListener('DOMContentLoaded', () => {
    document.body.classList.add('loaded');

    const savedTheme = localStorage.getItem('theme') || 'light';
    document.body.setAttribute("data-theme", savedTheme);
});

// see koodi lõik hoolitseb lehekülje dark/light tekkimise ja püsimise eest
document.getElementById("themeToggle").addEventListener("click", function() {
    let currentTheme = document.body.getAttribute("data-theme") === "dark" ? "light" : "dark";
    document.body.setAttribute("data-theme", currentTheme);
    localStorage.setItem('theme', currentTheme); // Salvesta teema eelistus
});