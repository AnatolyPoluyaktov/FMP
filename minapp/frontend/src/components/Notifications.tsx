import React, { useState, useEffect } from 'react';
import { Bell, Check, AlertTriangle, DollarSign, Calendar } from 'lucide-react';
import { apiService, Notification, NotificationStats } from '../services/api';

export const Notifications: React.FC = () => {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [stats, setStats] = useState<NotificationStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadNotifications();
  }, []);

  const loadNotifications = async () => {
    try {
      setLoading(true);
      const [notificationsData, statsData] = await Promise.all([
        apiService.getNotifications(),
        apiService.getNotificationStats()
      ]);
      setNotifications(notificationsData);
      setStats(statsData);
    } catch (error) {
      console.error('Error loading notifications:', error);
    } finally {
      setLoading(false);
    }
  };

  const markAsRead = async (id: string) => {
    try {
      await apiService.markNotificationAsRead(id);
      setNotifications(prev => 
        prev.map(notif => 
          notif.id === id ? { ...notif, is_read: true } : notif
        )
      );
      if (stats) {
        setStats(prev => prev ? { ...prev, unread_count: prev.unread_count - 1 } : null);
      }
    } catch (error) {
      console.error('Error marking notification as read:', error);
    }
  };

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case 'daily_reminder':
        return <Calendar className="notification-icon" />;
      case 'limit_warning':
        return <AlertTriangle className="notification-icon warning" />;
      case 'limit_exceeded':
        return <AlertTriangle className="notification-icon error" />;
      case 'income_reminder':
        return <DollarSign className="notification-icon" />;
      default:
        return <Bell className="notification-icon" />;
    }
  };

  const getNotificationClass = (type: string, isRead: boolean) => {
    let baseClass = 'notification-item';
    if (isRead) baseClass += ' read';
    
    switch (type) {
      case 'limit_warning':
        return baseClass + ' warning';
      case 'limit_exceeded':
        return baseClass + ' error';
      default:
        return baseClass;
    }
  };

  if (loading) {
    return (
      <div className="notifications-container">
        <div className="loading">Загрузка уведомлений...</div>
      </div>
    );
  }

  return (
    <div className="notifications-container">
      <div className="notifications-header">
        <h2>Уведомления</h2>
        {stats && stats.unread_count > 0 && (
          <div className="unread-badge">{stats.unread_count}</div>
        )}
      </div>

      {notifications.length === 0 ? (
        <div className="no-notifications">
          <Bell className="no-notifications-icon" />
          <p>Нет уведомлений</p>
        </div>
      ) : (
        <div className="notifications-list">
          {notifications.map(notification => (
            <div
              key={notification.id}
              className={getNotificationClass(notification.type, notification.is_read)}
            >
              <div className="notification-content">
                <div className="notification-header">
                  {getNotificationIcon(notification.type)}
                  <h3 className="notification-title">{notification.title}</h3>
                  {!notification.is_read && (
                    <button
                      className="mark-read-button"
                      onClick={() => markAsRead(notification.id)}
                      title="Отметить как прочитанное"
                    >
                      <Check size={16} />
                    </button>
                  )}
                </div>
                <p className="notification-message">{notification.message}</p>
                <div className="notification-time">
                  {new Date(notification.created_at).toLocaleString('ru-RU')}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
