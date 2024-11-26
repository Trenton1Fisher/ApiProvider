const apiUrl = 'https://api.dog.trentonfisher.xyz';

const tempButtonEvent = async () => {
    try {
        const res = await fetch(`${apiUrl}/getToken`)
        if (res.ok) {
            const data = await res.json();
            console.log(data);
        } else {
            console.error('Failed to fetch:', res.statusText);
        }
    } catch (error) {
        console.error('Error:', error);
    }
};

