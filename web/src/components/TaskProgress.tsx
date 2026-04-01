import { useEffect, useRef, useState } from 'react'
import { Modal, Progress, List, Typography, Badge } from 'antd'
import { useTranslation } from 'react-i18next'
import { useAuthStore } from '../store/auth'

const { Text } = Typography

interface TaskProgressProps {
  taskId: number | null
  open: boolean
  onClose: () => void
}

interface ProgressEvent {
  percent: number
  message: string
  level?: 'info' | 'error' | 'success'
}

export default function TaskProgress({ taskId, open, onClose }: TaskProgressProps) {
  const [percent, setPercent] = useState(0)
  const [logs, setLogs] = useState<ProgressEvent[]>([])
  const token = useAuthStore((s) => s.token)
  const esRef = useRef<EventSource | null>(null)
  const logsEndRef = useRef<HTMLDivElement>(null)
  const { t } = useTranslation()

  useEffect(() => {
    if (!open || taskId == null) return

    setPercent(0)
    setLogs([])

    const url = `/api/tasks/${taskId}/progress${token ? `?token=${encodeURIComponent(token)}` : ''}`
    const es = new EventSource(url)
    esRef.current = es

    es.onmessage = (e) => {
      try {
        const data: ProgressEvent = JSON.parse(e.data)
        setPercent(data.percent ?? 0)
        setLogs((prev) => [...prev, data])
      } catch {
        // ignore parse errors
      }
    }

    es.onerror = () => {
      es.close()
    }

    return () => {
      es.close()
      esRef.current = null
    }
  }, [open, taskId, token])

  useEffect(() => {
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [logs])

  function handleClose() {
    esRef.current?.close()
    esRef.current = null
    onClose()
  }

  const statusColor = percent >= 100 ? 'success' : 'active'

  return (
    <Modal
      title={t('tasks.taskProgressTitle', { id: taskId })}
      open={open}
      onCancel={handleClose}
      footer={null}
      width={600}
    >
      <Progress percent={percent} status={statusColor} style={{ marginBottom: 16 }} />
      <div
        style={{
          height: 240,
          overflowY: 'auto',
          border: '1px solid #f0f0f0',
          borderRadius: 4,
          padding: 8,
          background: '#fafafa',
        }}
      >
        <List
          size="small"
          dataSource={logs}
          renderItem={(item, idx) => (
            <List.Item key={idx} style={{ padding: '2px 0', border: 'none' }}>
              <Badge
                color={item.level === 'error' ? 'red' : item.level === 'success' ? 'green' : 'blue'}
                text={<Text style={{ fontSize: 12 }}>{item.message}</Text>}
              />
            </List.Item>
          )}
        />
        <div ref={logsEndRef} />
      </div>
    </Modal>
  )
}
