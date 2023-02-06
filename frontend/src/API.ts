
export const API = {
    async fetch(method: string, url: string, headers?: { [x: string]: string } | null, data?: any) {
        if (!headers) {
            headers = {'Content-Type': 'application/json'}
        }

        let token = localStorage.getItem('token')
        if (token !== null) {
            headers['Authorization'] = `Bearer ${token}`
        }

        const requestParameters: { [x: string]: any } = {
            method: method,
            headers: headers,
        }

        if (data) {
            requestParameters['body'] = JSON.stringify(data)
        }

        try {
            const resp = await fetch(url, requestParameters)
            if (!resp.ok) {
                console.log('Failed to remove token. Got status: ', resp.status)
            }

            if (resp.status === 401 || resp.status === 403) {
                localStorage.removeItem('token')
                window.location.href = '/'
            }

            if (!resp.ok) {
                throw new Error(resp.statusText);
            }

            return resp;
        } catch (error) {
            throw new Error('Something went wrong');
        }
    }
};
