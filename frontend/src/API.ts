
export const API = {
    async fetch(method: string, url: string, headers?: { [x: string]: string } | null, data?: any) {
        if (!headers) {
            headers = {'Content-Type': 'application/json'}
        }

        let token = localStorage.getItem('token')
        if (token !== null) {
            headers['Authorization'] = `Bearer ${token}`
        }
        const workspace_id = localStorage.getItem('workspace_id')
        if (workspace_id !== null) {
            headers['workspace-id'] = workspace_id
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
                localStorage.removeItem('workspace_id');
                localStorage.removeItem('workspace_name');
                window.location.href = '/'
            }

            if (!resp.ok) {
                return {}
            }

            const data = await resp.json()

            return data
        } catch (error) {
            return false
        }
    }
};
