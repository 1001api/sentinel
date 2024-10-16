const button = document.getElementById("g-auth");

button?.addEventListener("click", function() {
    const AUTH_URL = `${window.location.origin}/auth/google`;
    const width = 500;
    const height = 600;
    const left = (window.innerWidth / 2) - (width / 2);
    const top = (window.innerHeight / 2) - (height / 2);
    const popup = window.open(
        AUTH_URL,
        "Google Login",
        `width=${width},height=${height},top=${top},left=${left}`
    );

    const checkPopup = setInterval(() => {
        if (popup?.closed) clearInterval(checkPopup);
    }, 1000)
})

window.addEventListener("message", (e) => {
    // ensure retrieve only from origin
    if (e.origin !== window.location.origin) return;
    const { success } = e.data;
    if (success) window.location.href = `${window.location.origin}/events`;
})
