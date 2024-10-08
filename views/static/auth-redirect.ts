window.onload = () => {
    if (window.opener) {
        window.opener.postMessage({ success: true }, window.location.origin);
        window.close();
    }
}
