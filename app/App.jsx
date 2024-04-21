// cd ./app/
// npx expo start

import { StatusBar } from 'expo-status-bar';
import { StyleSheet, Text, View } from 'react-native';
import * as ImagePicker from 'expo-image-picker';
import { useState, useEffect } from 'react';

import Button from './Button';

import ImageViewer from './ImageViewer';
const PlaceholderImage = require('./assets/images/main-image.png')


//TODO: https://docs.expo.dev/tutorial/image-picker/#use-the-selected-image




export default function App() {

 // const [selectedImage, setSelectedImage] = useState(null);

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

  const [imageUrl, setImageUrl] = useState('');

  const updateImage = () => {
    fetch('http://192.168.233.170:8080/frame')
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.blob();
      })
      .then(blob => {
        //console.log(blob);
        // const uri = URL.createObjectURL(blob);
        // console.log("uri: ", uri);
        // setImageUrl(uri);
        const reader = new FileReader();
        reader.onload = () => {
          const dataUrl = reader.result;
          setImageUrl(dataUrl);
        };
        reader.readAsDataURL(blob);
      })
      .catch(error => {
        console.error('Error fetching image:', error);
      });
  };

  useEffect(() => {
    //const interval = setInterval(updateImage, 2000);
    updateImage(); // Initial call

    // Clean up the interval
   //return () => clearInterval(interval);
  }, []);
   

  return (
    <View style={styles.container}>
      <Text style={styles.heading}>
        Live Shapes
      </Text>
      {/* <ImageViewer 
        placeholderImageSource={PlaceholderImage}
        selectedImage={selectedImage} /> */}
      {/* <ImageViewer
        placeholderImageSource={{ uri: imageUrl }}
        selectedImage={{ uri: imageUrl }}
        style={{ width: 100, height: 100 }} 
      /> */}

      

      <Button label="Upload" onPress={pickImageAsync} />
      <Button label="Take a Photo" onPress={pickImageAsync} />
      <StatusBar style="auto"/>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: 'black',
    color: 'black',
    alignItems: 'center',
    justifyContent: 'center',
    //gap: 25,
  },
  
  heading: {
    color: 'white',
    fontSize: 36,
    //paddingVertical: 17,
    fontFamily: 'monospace',
    fontWeight: 'bold',
    //marginBottom: 20,
  },

  gif: {
    //aspectRatio: 1,
    width: 300,
    height: 100,
  }

});
