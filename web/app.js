const apiUrl = 'https://api.dog.trentonfisher.xyz';

const tempButtonEvent = async () => {
    try {
        const res = await fetch(`${apiUrl}/api/get-token`)
        if (res.ok) {
            const data = await res.json();
            console.log(data);
        } else {
            console.error('Failed to fetch:', res.statusText);
        }
    } catch (error) {
        console.log('Error:', error);
    }
};

