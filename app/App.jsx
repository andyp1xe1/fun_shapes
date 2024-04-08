// cd ./app/
// npx expo start

import { StatusBar } from 'expo-status-bar';
import { StyleSheet, Text, View } from 'react-native';
import * as ImagePicker from 'expo-image-picker';
import { useState } from 'react';
import {encode} from 'base-64';

import Button from './Button';

import ImageViewer from './ImageViewer';
const PlaceholderImage = require('./assets/images/main-image.png')


//TODO: https://docs.expo.dev/tutorial/image-picker/#use-the-selected-image




export default function App() {

 






  const [selectedImage, setSelectedImage] = useState(null);

  const pickImageAsync = async () => {
    let result = await ImagePicker.launchImageLibraryAsync({
      allowsEditing: true,
      quality: 1,
    });

    if (!result.canceled) {
      console.log(result);
      setSelectedImage(result.assets[0].uri);
      
    } else {
      alert('You did not select any image.');
    }
  };



   // websockets implementation
   const ws = new WebSocket('ws://192.168.233.170:8080/ws');
   ws.onopen = () => {
    const base64 = `data:${selectedImage.type};base64,${encode(
      selectedImage.data
    )}`;
    ws.send(base64);
    console.log(base64);
    console.log(selectedImage);
  };
   
   ws.onmessage = e => {
     // a message was received
     console.log(e.data);
   };
   ws.onerror = e => {
     // an error occurred
     console.log(e.message);
   };
   ws.onclose = e => {
     // connection closed
     console.log(e.code, e.reason);
   };
   

  return (
    <View style={styles.container}>
      <Text style={styles.heading}>
        Fun Shapes
      </Text>
      <ImageViewer 
        placeholderImageSource={PlaceholderImage}
        selectedImage={selectedImage} />
      <Button label="Upload" onPress={pickImageAsync} />
      <Button label="Take a Photo" onPress={pickImageAsync} />
      <StatusBar style="auto"/>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: 'white',
    color: 'black',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 25,
  },
  
  heading: {
    color: '#4d676e',
    fontSize: 36,
    paddingVertical: 16,
    fontFamily: 'monospace',
    fontWeight: 'bold',
  },
});
