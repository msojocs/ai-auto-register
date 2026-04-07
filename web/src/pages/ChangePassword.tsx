import { useState } from 'react'
import { Card, Form, Input, Button, message, Space, Typography } from 'antd'
import { useTranslation } from 'react-i18next'
import { changePassword } from '../api/auth'

const { Title } = Typography

export default function ChangePassword() {
  const { t } = useTranslation()
  const [form] = Form.useForm<{ current_password: string; new_password: string; confirm_password: string }>()
  const [saving, setSaving] = useState(false)

  async function handleSubmit(values: { current_password: string; new_password: string; confirm_password: string }) {
    setSaving(true)
    try {
      await changePassword(values.current_password, values.new_password)
      message.success(t('userMenu.passwordChanged'))
      form.resetFields()
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : t('userMenu.passwordChangeFailed')
      message.error(errorMsg)
    } finally {
      setSaving(false)
    }
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <Title level={3} style={{ marginBottom: 0 }}>
        {t('userMenu.changePassword')}
      </Title>
      <Card style={{ maxWidth: 560 }}>
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="current_password"
            label={t('userMenu.currentPassword')}
            rules={[{ required: true, message: t('userMenu.currentPasswordRequired') }]}
          >
            <Input.Password autoComplete="current-password" />
          </Form.Item>
          <Form.Item
            name="new_password"
            label={t('userMenu.newPassword')}
            rules={[
              { required: true, message: t('userMenu.newPasswordRequired') },
              { min: 6, message: t('userMenu.newPasswordMinLength') },
            ]}
          >
            <Input.Password autoComplete="new-password" />
          </Form.Item>
          <Form.Item
            name="confirm_password"
            label={t('userMenu.confirmPassword')}
            dependencies={['new_password']}
            rules={[
              { required: true, message: t('userMenu.confirmPasswordRequired') },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('new_password') === value) {
                    return Promise.resolve()
                  }
                  return Promise.reject(new Error(t('userMenu.passwordNotMatch')))
                },
              }),
            ]}
          >
            <Input.Password autoComplete="new-password" />
          </Form.Item>
          <Form.Item style={{ marginBottom: 0 }}>
            <Button type="primary" htmlType="submit" loading={saving}>
              {t('common.save')}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </Space>
  )
}
