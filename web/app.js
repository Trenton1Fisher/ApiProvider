const apiUrl = 'https://api.dog.trentonfisher.xyz';

const tokenDisplay = document.getElementById('token-display');

const tempButtonEvent = async () => {
    try {
        const res = await fetch(`${apiUrl}/api/get-token`);
        if (res.ok) {
            const data = await res.json();
            if (data) {
                const token = data["token"];
                // Display the token inside a div with an input element for easy copying
                tokenDisplay.innerHTML = `
                    <div id="api-token-container">
                        <p>Token:</p>
                        <input type="text" id="api-token" value="${token}" readonly />
                    </div>
                `;
            }
        } else {
            console.error('Failed to fetch:', res.statusText);
        }
    } catch (error) {
        console.log('Error:', error);
    }
};

const copyToken = () => {
    const tokenInput = document.getElementById('api-token');
    if (tokenInput) {
        tokenInput.select();
        tokenInput.setSelectionRange(0, 99999); // For mobile devices

        navigator.clipboard.writeText(tokenInput.value).then(() => {
            alert('Token copied to clipboard: ' + tokenInput.value);
        }).catch(err => {
            console.error('Failed to copy token: ', err);
        });
    }
};
