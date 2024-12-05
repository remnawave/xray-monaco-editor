import { fetchWithProgress } from '@/shared/utils/fetch-with-progress'
import { useEffect, useState } from 'react'

import { LoadingScreen } from '@/shared/ui/loading-screen'
import { DEFAULT_CONFIG } from '@/shared/constants'

import { ConfigPageComponent } from '../components/config.page.component'

export function ConfigPageConnector() {
    const [downloadProgress, setDownloadProgress] = useState(0)
    const [isLoading, setIsLoading] = useState(true)

    const config = DEFAULT_CONFIG
    useEffect(() => {
        const initWasm = async () => {
            try {
                const go = new window.Go()
                const wasmInitialized = new Promise<void>((resolve) => {
                    window.onWasmInitialized = () => {
                        console.info('WASM module initialized')
                        resolve()
                    }
                })

                const wasmBytes = await fetchWithProgress('main.wasm', setDownloadProgress)

                const { instance } = await WebAssembly.instantiate(wasmBytes, go.importObject)
                go.run(instance)
                await wasmInitialized

                if (typeof window.XrayParseConfig === 'function') {
                    setIsLoading(false)
                } else {
                    throw new Error('XrayParseConfig not initialized')
                }
            } catch (err: unknown) {
                console.error('WASM initialization error:', err)
                setIsLoading(false)
            }
        }

        initWasm()
        return () => {
            delete window.onWasmInitialized
        }
    }, [])

    if (isLoading) {
        return <LoadingScreen text={`WASM module is loading...`} value={downloadProgress} />
    }

    return <ConfigPageComponent config={config} />
}
