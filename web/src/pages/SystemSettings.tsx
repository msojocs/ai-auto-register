import { Tabs } from 'antd'
import { useTranslation } from 'react-i18next'
import ProxyGroupManager from '../components/ProxyGroupManager'

export default function SystemSettings() {
  const { t } = useTranslation()

  return (
    <Tabs
      items={[
        {
          key: 'proxy-groups',
          label: t('proxyGroups.title'),
          children: <ProxyGroupManager />,
        },
      ]}
    />
  )
}