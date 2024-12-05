import { Container } from '@mantine/core'

import { ConfigEditorWidget } from '@widgets/config/config-editor'
import { HeaderWidget } from '@widgets/header'
import { Page } from '@/shared/ui/page'

import { Props } from './interfaces'

export const ConfigPageComponent = (props: Props) => {
    const { config } = props

    return (
        <>
            <HeaderWidget />

            <Container size="lg">
                <Page title="Config">
                    <ConfigEditorWidget config={config} />
                </Page>
            </Container>
        </>
    )
}
