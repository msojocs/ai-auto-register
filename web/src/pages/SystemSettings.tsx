import { useEffect, useState } from 'react'
import { Tabs, Form, Input, Button, Card, Space, Typography, message, Select } from 'antd'
import { useTranslation } from 'react-i18next'
import ProxyGroupManager from '../components/ProxyGroupManager'
import { getSystemSettings, updateSystemSettings } from '../api/settings'
import { getProxyGroups, type ProxyGroup } from '../api/proxyGroups'

const { Text } = Typography

export default function SystemSettings() {
  const { t } = useTranslation()
  const [form] = Form.useForm<{
    sentinel_base_url: string
    account_action_proxy_group_id?: number
  }>()
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [proxyGroups, setProxyGroups] = useState<ProxyGroup[]>([])

  async function fetchSettings() {
    setLoading(true)
    try {
      const { data } = await getSystemSettings()
      form.setFieldsValue({
        sentinel_base_url: data.sentinel_base_url,
        account_action_proxy_group_id: data.account_action_proxy_group_id,
      })
    } catch {
      message.error(t('settings.failedToLoad'))
    } finally {
      setLoading(false)
    }
  }

  async function fetchProxyGroups() {
    try {
      const { data } = await getProxyGroups()
      setProxyGroups(data.groups ?? [])
    } catch {
      message.error(t('proxyGroups.failedToLoad'))
    }
  }

  useEffect(() => {
    void fetchSettings()
    void fetchProxyGroups()
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  async function handleSave(values: {
    sentinel_base_url: string
    account_action_proxy_group_id?: number
  }) {
    setSaving(true)
    try {
      const { data } = await updateSystemSettings(values)
      form.setFieldsValue({
        sentinel_base_url: data.sentinel_base_url,
        account_action_proxy_group_id: data.account_action_proxy_group_id,
      })
      message.success(t('settings.saved'))
    } catch {
      message.error(t('settings.failedToSave'))
    } finally {
      setSaving(false)
    }
  }

  return (
    <Tabs
      items={[
        {
          key: 'runtime',
          label: t('settings.runtimeTab'),
          children: (
            <Card loading={loading}>
              <Form form={form} layout="vertical" onFinish={handleSave}>
                <Form.Item
                  label={t('settings.sentinelBaseUrl')}
                  name="sentinel_base_url"
                  rules={[{ required: true, message: t('settings.sentinelBaseUrlRequired') }]}
                >
                  <Input placeholder={t('settings.sentinelBaseUrlPlaceholder')} />
                </Form.Item>
                <Form.Item label={t('settings.accountActionProxyGroup')} name="account_action_proxy_group_id">
                  <Select
                    allowClear
                    options={proxyGroups.map((group) => ({ label: group.name, value: group.id }))}
                    placeholder={t('settings.proxyGroupPlaceholder')}
                  />
                </Form.Item>
                <Space direction="vertical" size={12} style={{ width: '100%' }}>
                  <Text type="secondary">{t('settings.sentinelBaseUrlHelp')}</Text>
                  <Text type="secondary">{t('settings.proxyGroupHelp')}</Text>
                  <Button type="primary" htmlType="submit" loading={saving}>
                    {t('common.save')}
                  </Button>
                </Space>
              </Form>
            </Card>
          ),
        },
        {
          key: 'proxy-groups',
          label: t('proxyGroups.title'),
          children: <ProxyGroupManager />,
        },
      ]}
    />
  )
}
