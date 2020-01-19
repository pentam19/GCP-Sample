import base64
import requests

# Cloud Functions Pub/Sub Trigger
def hello_pubsub(event, context):
    """Triggered from a message on a Cloud Pub/Sub topic.
    Args:
         event (dict): Event payload.
         context (google.cloud.functions.Context): Metadata for the event.
    """
    pubsub_message = base64.b64decode(event['data']).decode('utf-8')
    print(pubsub_message)
    
    headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36'}
    url = 'https://www.google.com/search?q=%22{}%22'.format(pubsub_message)
    print(url)
    #response = requests.get(url,  headers=headers)
    #print(response.status_code)    # HTTPのステータスコード取得
    #print(response.text)    # レスポンスのHTMLを文字列で取得
    # POST
    """
    response = requests.post('http://www.example.com', data={'foo': 'bar'})
    print(response.status_code)    # HTTPのステータスコード取得
    print(response.text)    # レスポンスのHTMLを文字列で取得
    """