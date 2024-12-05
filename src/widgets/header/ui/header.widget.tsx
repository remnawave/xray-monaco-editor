import { ActionIcon, Group, Title } from '@mantine/core'
import { StickyHeader } from '@/shared/ui/sticky-header'
import { PiGithubLogo } from 'react-icons/pi'

import classes from './header.module.css'

export function HeaderWidget() {
    return (
        <StickyHeader className={classes.root} px="md">
            <Group h="100%" justify="space-between">
                <Title order={3}>Xray Config Validator</Title>

                <ActionIcon
                    component="a"
                    href="https://github.com/remnawave/xray-monaco-editor"
                    size="lg"
                    target="_blank"
                    variant="subtle"
                >
                    <PiGithubLogo size={24} />
                </ActionIcon>
            </Group>
        </StickyHeader>
    )
}
