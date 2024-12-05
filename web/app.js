const apiUrl = 'https://api.dog.trentonfisher.xyz';

const tokenDisplay = document.getElementById('token-display')

const tempButtonEvent = async () => {
    try {
        const res = await fetch(`${apiUrl}/api/get-token`)
        if (res.ok) {
            const data = await res.json();
            if(data){
                console.log(data)
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
    tokenInput.select();
    tokenInput.setSelectionRange(0, 99999);

    navigator.clipboard.writeText(tokenInput.value).then(() => {
      alert('Token copied to clipboard: ' + tokenInput.value);
    }).catch(err => {
      console.error('Failed to copy token: ', err);
    });
  };