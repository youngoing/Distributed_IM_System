// components/EmojiPicker.js

import data from '@emoji-mart/data';
import Picker from '@emoji-mart/react';

const EmojiPicker = ({ onEmojiSelect }) => {
  return (
    <div className="emoji-picker">
      <Picker data={data} onEmojiSelect={onEmojiSelect} />
    </div>
  );
};

export default EmojiPicker;