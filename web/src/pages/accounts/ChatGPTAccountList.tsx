import type { ColumnsType } from 'antd/es/table'
import { useCallback, useMemo, useState } from 'react'
import { Button, Modal, Space, message } from 'antd'
import { ReloadOutlined, EyeOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import {
  getChatGPTAccountDetail,
  refreshChatGPTToken,
  type Account,
  type ChatGPTAccountDetailResult,
} from '../../api/accounts'
import AccountTableTemplate from './AccountTableTemplate'
import { parseAccountExtra, extractChatGPTAccountIdFromToken, extractJwtExp } from './accountExtra'

export default function ChatGPTAccountList() {
  const { t } = useTranslation()
  const [refreshingId, setRefreshingId] = useState<number | null>(null)
  const [detailLoadingId, setDetailLoadingId] = useState<number | null>(null)
  const [detailOpen, setDetailOpen] = useState(false)
  const [detail, setDetail] = useState<ChatGPTAccountDetailResult | null>(null)

  const columns: ColumnsType<Account> = useMemo(() => [
    {
      title: t('accounts.accountId'),
      key: 'account_id',
      render: (_, record) => {
        const extra = parseAccountExtra(record.extra)
        const fromExtra = extra.account_id
        if (typeof fromExtra === 'string' && fromExtra) return fromExtra

        const accessToken = extra.access_token
        if (typeof accessToken !== 'string') return '-'
        const parsed = extractChatGPTAccountIdFromToken(accessToken)
        return parsed || '-'
      },
    },
    {
      title: t('accounts.accessTokenExpiresAt'),
      key: 'access_token_expires_at',
      render: (_, record) => {
        const extra = parseAccountExtra(record.extra)
        const accessToken = extra.access_token
        if (typeof accessToken !== 'string') return '-'
        const exp = extractJwtExp(accessToken)
        if (!exp) return '-'
        return new Date(exp * 1000).toLocaleString()
      },
    },
  ], [t])

  const handleRefreshToken = useCallback(async (record: Account) => {
    setRefreshingId(record.id)
    try {
      await refreshChatGPTToken(record.id)
      message.success(t('accounts.refreshTokenSuccess'))
    } catch {
      message.error(t('accounts.refreshTokenFailed'))
    } finally {
      setRefreshingId(null)
    }
  }, [t])

  const handleViewDetail = useCallback(async (record: Account) => {
    setDetailOpen(true)
    setDetailLoadingId(record.id)
    try {
      const { data } = await getChatGPTAccountDetail(record.id)
      setDetail(data)
    } catch (err) {
      message.error(
        t('accounts.detailFailed', {
          message: err instanceof Error ? err.message : 'unknown error',
        }),
      )
      setDetail(null)
    } finally {
      setDetailLoadingId(null)
    }
  }, [t])

  const renderExtraActions = useCallback((record: Account) => (
    <Space size={4}>
      <Button
        size="small"
        icon={<ReloadOutlined />}
        loading={refreshingId === record.id}
        onClick={() => void handleRefreshToken(record)}
      >
        {t('accounts.refreshToken')}
      </Button>
      <Button
        size="small"
        icon={<EyeOutlined />}
        loading={detailLoadingId === record.id}
        onClick={() => void handleViewDetail(record)}
      >
        {t('accounts.detail')}
      </Button>
    </Space>
  ), [detailLoadingId, handleRefreshToken, handleViewDetail, refreshingId, t])

  return (
    <>
      <AccountTableTemplate
        title={t('accounts.chatgptTitle')}
        accountType="chatgpt"
        extraColumns={columns}
        hideTypeColumn
        renderExtraActions={renderExtraActions}
      />
      <Modal
        title={t('accounts.detailTitle')}
        open={detailOpen}
        onCancel={() => setDetailOpen(false)}
        footer={null}
        width={820}
      >
        <pre style={{ maxHeight: '60vh', overflow: 'auto', margin: 0 }}>
          {detailLoadingId !== null
            ? t('common.loading')
            : detail
              ? JSON.stringify(detail, null, 2)
              : t('accounts.detailEmpty')}
        </pre>
      </Modal>
    </>
  )
}
