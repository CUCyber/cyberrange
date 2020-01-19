(function() {
  var conn = undefined;

  var chat = {
    messageToSend: '',
    init: function() {
      this.cacheDOM();
      this.bindEvents();
    },
    cacheDOM: function() {
      this.chatHistory = $('.chat-history');
      this.button = $('#send-button');
      this.textarea = $('#message-to-send');
      this.chatHistoryList = this.chatHistory.find('ul');
    },
    bindEvents: function() {
      this.button.on('click', this.sendMessage.bind(this));
      this.textarea.on('keyup', this.sendMessageEnter.bind(this));
    },
    render: function(message) {
      if (message.Data != '') {
          this.scrollToBottom();
          var messageElement = document.createElement('div');

          messageElement.className = "message";
          messageElement.innerText = `[${message.Time}] ${message.Name}: ${message.Data}`;

          this.chatHistoryList.append(messageElement);
          this.scrollToBottom();
          this.textarea.val('');
      }
    },
    sendMessage: function() {
      this.messageToSend = this.textarea.val()
      if (conn && this.messageToSend.trim() !== '') {
        conn.send(JSON.stringify({
          "Data": this.messageToSend
        }));
      }
    },
    sendMessageEnter: function(event) {
      if (event.keyCode === 13) {
        this.sendMessage();
      }
    },
    scrollToBottom: function() {
      this.chatHistory.scrollTop(this.chatHistory[0].scrollHeight);
    },
  };

  if (window["WebSocket"]) {
    conn = new WebSocket("ws://" + document.location.host + "/chat/ws");

    conn.onmessage = function(event) {
      const json = (function(resp) {
        try {
          return JSON.parse(resp);
        } catch (err) {
          return false;
        }
      })(event.data);

      if (!json) {
        toastr.error('Action Request Failed.', 'Chat');
        return;
      }

      if (json.Messages) {
        json.Messages.forEach(el => {
          chat.render(el);
        });
      } else {
        chat.render(json);
      }
    };

    $(window).on('beforeunload', function() {
      conn.close();
    });

    chat.init();
  }
})();
