import axios from 'axios'

export const fetchWithProgress = async (url: string, onProgress?: (progress: number) => void) => {
    try {
        const response = await axios.get(url, {
            onDownloadProgress: (progressEvent) => {
                if (progressEvent.total) {
                    const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
                    onProgress?.(progress)
                } else {
                    onProgress?.(100)
                }
            },
            responseType: 'arraybuffer'
        })

        return response.data
    } catch (error) {
        console.error('Download failed:', error)
        throw error
    }
}
