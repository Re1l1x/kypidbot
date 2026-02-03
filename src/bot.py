#!/usr/bin/env python3

import logging

from telegram import InlineKeyboardButton, InlineKeyboardMarkup, Update
from telegram.ext import Application, CallbackQueryHandler, CommandHandler, ContextTypes

from config import TELEGRAM_TOKEN
from database import Database

logging.basicConfig(
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    level=logging.INFO,
)
logger = logging.getLogger(__name__)
logging.getLogger("httpx").setLevel(logging.WARNING)

# Initialize database
db = Database()


async def start(update: Update, _: ContextTypes.DEFAULT_TYPE) -> None:
    if not update.message or not update.effective_user:
        return

    user = update.effective_user

    # Save user to database
    db.save_user(
        telegram_id=user.id,
        username=user.username,
        first_name=user.first_name,
        last_name=user.last_name,
        is_bot=user.is_bot,
        language_code=user.language_code,
        is_premium=user.is_premium,
    )

    await update.message.reply_text(
        f"👋 Привет, {user.first_name}!\n\n"
        f"Я бот для поиска пар на День Святого Валентина! 💕\n\n"
        f"Расскажи мне о своих интересах, и я помогу найти тебе пару с похожими увлечениями."
    )


TIME_RANGES = [
    "10:00 -- 12:00",
    "12:00 -- 14:00",
    "14:00 -- 16:00",
    "16:00 -- 18:00",
    "18:00 -- 20:00",
    "20:00 -- 22:00",
]


def binary_to_set(binary_str: str) -> set[str]:
    """Convert binary string to set of selected time ranges.

    Example: '101000' -> {'10:00 -- 12:00', '14:00 -- 16:00'}
    """
    selected = set()
    for i, bit in enumerate(binary_str):
        if bit == "1" and i < len(TIME_RANGES):
            selected.add(TIME_RANGES[i])
    return selected


def set_to_binary(selected: set[str]) -> str:
    """Convert set of selected time ranges to binary string.

    Example: {'10:00 -- 12:00', '14:00 -- 16:00'} -> '101000'
    """
    binary = []
    for time_range in TIME_RANGES:
        binary.append("1" if time_range in selected else "0")
    return "".join(binary)


def create_time_keyboard(selected_times: set[str]) -> InlineKeyboardMarkup:
    """Create keyboard with checkmarks for selected time ranges."""
    keyboard = []
    for i in range(0, len(TIME_RANGES), 2):
        row = []
        for time_range in TIME_RANGES[i : i + 2]:
            text = f"> {time_range} <" if time_range in selected_times else time_range
            row.append(InlineKeyboardButton(text, callback_data=f"time_{time_range}"))
        keyboard.append(row)

    return InlineKeyboardMarkup(keyboard)


async def time_command(update: Update, _: ContextTypes.DEFAULT_TYPE) -> None:
    if not update.message or not update.effective_user:
        return

    user = update.effective_user

    # Ensure user exists in database
    db.save_user(
        telegram_id=user.id,
        username=user.username,
        first_name=user.first_name,
        last_name=user.last_name,
        is_bot=user.is_bot,
        language_code=user.language_code,
        is_premium=user.is_premium,
    )

    # Load time ranges from database
    binary_str = db.get_time_ranges(user.id)
    selected_times = binary_to_set(binary_str)

    keyboard = create_time_keyboard(selected_times)

    await update.message.reply_text(
        "Выберите удобные промежутки, в которые вам будет зарандомлена свиданка.\n\nP.S. свиданка будет не позже, чем конец промежутка минус пол часа",
        reply_markup=keyboard,
    )


async def time_button_callback(update: Update, _: ContextTypes.DEFAULT_TYPE) -> None:
    """Handle time range button clicks."""
    query = update.callback_query
    if not query or not query.data or not query.from_user:
        return

    await query.answer()

    user = query.from_user

    # Extract time range from callback data
    time_range = query.data.replace("time_", "")

    # Load current selection from database
    binary_str = db.get_time_ranges(user.id)
    selected = binary_to_set(binary_str)

    # Toggle selection
    if time_range in selected:
        selected.remove(time_range)
    else:
        selected.add(time_range)

    # Save to database
    new_binary = set_to_binary(selected)
    db.save_time_ranges(user.id, new_binary)

    # Update keyboard
    keyboard = create_time_keyboard(selected)

    await query.edit_message_reply_markup(reply_markup=keyboard)


def main() -> None:
    application = Application.builder().token(TELEGRAM_TOKEN).build()

    # Register command handlers
    application.add_handler(CommandHandler("start", start))
    application.add_handler(CommandHandler("time", time_command))

    # Register callback query handler for time selection
    application.add_handler(
        CallbackQueryHandler(time_button_callback, pattern="^time_")
    )

    logger.info("Starting bot...")
    application.run_polling(allowed_updates=Update.ALL_TYPES)


if __name__ == "__main__":
    main()
