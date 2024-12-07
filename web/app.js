const apiUrl = 'https://api.dog.trentonfisher.xyz';

const tokenDisplay = document.getElementById('token-display');

const tempButtonEvent = async () => {
    try {
        const res = await fetch(`${apiUrl}/api/get-token`);
        if (res.ok) {
            const data = await res.json();
            if (data) {
                const token = data["token"];
                tokenDisplay.innerHTML = `
                    <div id="api-token-container">
                        <p>Token:</p>
                        <p classname="token-copy" id="api-token">${token}</p>
                        <p>Please Do not Share With anyone</p>
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
