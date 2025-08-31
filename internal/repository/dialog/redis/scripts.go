package redis

// LUA скрипты для работы с диалогами в Redis

const (
	// SendMessageScript - скрипт для отправки сообщения и увеличения счетчика диалогов
	SendMessageScript = `
		-- Получаем ключи и аргументы
		local dialog_key = KEYS[1]
		local counter_key = KEYS[2]
		local message_json = ARGV[1]
		local ttl_seconds = tonumber(ARGV[2])
		
		-- Добавляем сообщение в список диалога
		local result = redis.call('LPUSH', dialog_key, message_json)
		
		-- Устанавливаем TTL для ключа диалога
		redis.call('EXPIRE', dialog_key, ttl_seconds)
		
		-- Увеличиваем счетчик диалогов
		local counter = redis.call('INCR', counter_key)
		
		-- Устанавливаем TTL для счетчика (больше чем для диалогов)
		redis.call('EXPIRE', counter_key, ttl_seconds * 2)
		
		-- Возвращаем количество сообщений в диалоге и значение счетчика
		return {result, counter}
	`

	// GetDialogCountScript - скрипт для получения количества диалогов
	GetDialogCountScript = `
		-- Получаем ключ счетчика
		local counter_key = KEYS[1]
		
		-- Получаем значение счетчика
		local counter = redis.call('GET', counter_key)
		
		-- Если счетчик не существует, возвращаем 0
		if counter == false then
			return 0
		end
		
		return tonumber(counter)
	`

	// GetDialogStatsScript - скрипт для получения статистики диалогов
	GetDialogStatsScript = `
		-- Получаем ключи
		local counter_key = KEYS[1]
		local pattern = ARGV[1]
		
		-- Получаем количество диалогов
		local total_dialogs = redis.call('GET', counter_key)
		if total_dialogs == false then
			total_dialogs = 0
		else
			total_dialogs = tonumber(total_dialogs)
		end
		
		-- Получаем количество активных диалогов (ключей, соответствующих паттерну)
		local active_dialogs = 0
		local keys = redis.call('KEYS', pattern)
		if keys then
			active_dialogs = #keys
		end
		
		-- Возвращаем статистику
		return {total_dialogs, active_dialogs}
	`
)
